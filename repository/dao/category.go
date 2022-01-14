package dao

import (
	"app/lib/nestedset"
	"app/util"
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

type Category struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" nestedset:"id" json:"id"`
	Name          string         `gorm:"size:200;uniqueIndex;not null" json:"name"`
	Description   string         `gorm:"type:text" json:"description"`
	Amount        uint           `gorm:"default:0" binding:"-" json:"amount"`
	Posts         []Post         `gorm:"foreignkey:categoryID" binding:"-" json:"posts"`
	UserID        string         `json:"userID"`
	User          *User          `gorm:"foreignkey:UserID" binding:"-" json:"user,omitempty"`
	ParentID      sql.NullInt64  `nestedset:"parent_id" json:"-"`
	Parent        *Category      `gorm:"foreignkey:ParentID" binding:"-" json:"parent"`
	Rgt           int            `nestedset:"rgt" json:"left"`
	Lft           int            `nestedset:"lft" json:"right"`
	Depth         int            `nestedset:"depth" json:"depth"`
	ChildrenCount int            `nestedset:"children_count" json:"childrenCount"`
	Children      []Category     `gorm:"-" binding:"-" json:"children"`
	Parents       []Category     `gorm:"-" binding:"-" json:"parents"`
	CreatedAt     util.LocalTime `json:"createdAt"`
	UpdatedAt     util.LocalTime `json:"updatedAt"`
	DeletedAt     util.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (m Category) Create(parent *Category) (Category, error) {
	m.ParentID = sql.NullInt64{Valid: true, Int64: parent.ID}
	if err := nestedset.Create(db, &m, parent); err != nil {
		return m, err
	}
	return m, nil
}

func (m Category) MoveTo(parent *Category) error {
	return nestedset.MoveTo(db, m, parent, nestedset.MoveDirectionInner)
}

func (m Category) Update(values interface{}) (Category, error) {
	err := db.Model(&m).Updates(values).Error
	return m, err
}

func (m Category) Save(cols []string) (Category, error) {
	tx := db.Session(&gorm.Session{})
	if len(cols) > 0 {
		for _, col := range cols {
			tx = tx.Select(col)
		}
	}
	if err := tx.Updates(&m).Error; err != nil {
		return m, err
	}
	return m, nil
}

func (m Category) ancestor() ([]Category, error) {
	var rows []Category
	if err := db.Model(&Category{}).Where("lft < ? AND rgt > ? AND parent_id IS NOT NULL", m.Lft, m.Rgt).Order("lft asc").Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func (m Category) descendant() ([]Category, error) {
	var rows []Category
	if err := db.Model(&Category{}).Where("lft > ? AND rgt < ?", m.Lft, m.Rgt).Order("lft asc").Find(&rows).Error; err != nil {
		return rows, err
	}
	descendants := categoryDescendantTree(rows, &m)
	return descendants, nil
}

func categoryDescendantTree(rows []Category, parent *Category) []Category {
	next := make([]Category, 0)
	for i := 0; i < len(rows); i++ {
		if rows[i].ParentID.Valid && rows[i].ParentID.Int64 == parent.ID {
			rows[i].Children = categoryDescendantTree(rows, &rows[i])
			next = append(next, rows[i])
		}
	}
	return next
}

func FindCategory(id uint, options map[string]interface{}) (Category, error) {
	var one Category
	if err := db.Scopes(applyQueryOptions(options)).First(&one, "id = ?", id).Error; err != nil {
		return one, err
	}
	return one, nil
}

func FindCategoryHierarchy(id uint, options map[string]interface{}) (Category, error) {
	one, err := FindCategory(id, options)
	if err != nil {
		return one, err
	}
	// if !one.ParentID.Valid {
	// 	return one, errors.New("不可访问根文件夹")
	// }
	descendants, err := one.descendant()
	if err != nil {
		return one, err
	}
	one.Children = descendants
	ancestors, err := one.ancestor()
	if err != nil {
		return one, err
	}
	one.Parents = ancestors
	return one, nil
}

func CategoryExists(id uint) (bool, Category) {
	var one Category
	err := db.Where("id = ?", id).First(&one).Error
	notFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !notFound, one
}

func FindCategories(options map[string]interface{}) ([]Category, error) {
	var rows []Category
	if err := db.Where("parent_id IS NOT NULL").Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, err
	}
	return rows, nil
}

func FindAndCountCategories(options map[string]interface{}) ([]Category, int64, error) {
	var rows []Category
	var count int64
	if err := db.Where("parent_id IS NOT NULL").Scopes(applyQueryOptions(options)).Find(&rows).Error; err != nil {
		return rows, count, err
	}
	delete(options, "offset")
	delete(options, "limit")
	if err := db.Model(&Category{}).Where("parent_id IS NOT NULL").Scopes(applyQueryOptions(options)).Count(&count).Error; err != nil {
		return rows, count, err
	}
	return rows, count, nil
}

func (m Category) Delete() error {
	tx := db.Begin()
	if !m.ParentID.Valid {
		return errors.New("不可删除根文件夹")
	}
	err := tx.Model(&m).Association("Posts").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	var children []Category
	err = tx.Model(&Category{}).Where("parent_id = ?", m.ID).Find(&children).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var root Category
	err = tx.First(&root, 1).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, child := range children {
		err := child.MoveTo(&root)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Delete(&m).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func DeleteCategories(ids []uint) error {
	tx := db.Begin()
	for _, id := range ids {
		if id == 1 {
			return errors.New("不可删除根文件夹")
		}
	}
	var rows []Category
	err := tx.Where("id IN (?)", ids).Find(&rows).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(rows) > 0 {
		depth := rows[0].Depth
		for _, row := range rows {
			if row.Depth != depth {
				return errors.New("只能批量删除同级文件夹")
			}
		}
	}
	err = tx.Model(&rows).Association("Posts").Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	var root Category
	err = tx.First(&root, 1).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var children []Category
	err = tx.Model(&Category{}).Where("parent_id IN (?)", ids).Find(&children).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, child := range children {
		err := child.MoveTo(&root)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Delete(&rows).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m Category) Relations(col string) *gorm.Association {
	return db.Model(&m).Association(col)
}

func (m *Category) Add(next []Post) (err error) {
	tx := db.Begin()
	err = tx.Model(&m).Select("amount").Updates(map[string]interface{}{"amount": gorm.Expr("amount + ?", len(next))}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	m.Amount += uint(len(next))
	if len(next) > 0 {
		err = tx.Model(&m).Association("Posts").Append(next)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return
}

func (m *Category) Remove(next []Post) (err error) {
	tx := db.Begin()
	if len(next) > 0 {
		err = tx.Model(&m).Select("amount").Updates(map[string]interface{}{"amount": gorm.Expr("amount - ?", len(next))}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		m.Amount -= uint(len(next))
		if err = tx.Model(&m).Association("Posts").Delete(next); err != nil {
			tx.Rollback()
			return err
		}
	} else if len(next) == 0 {
		if err = tx.Model(&m).Association("Posts").Clear(); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return
}

func MoveCategory(from *Category, to *Category, fromRows []Post, toRows []Post) (err error) {
	tx := db.Begin()
	err = tx.Model(from).Select("amount").Updates(map[string]interface{}{"amount": gorm.Expr("amount - ?", len(fromRows))}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(fromRows) == 0 {
		if err = tx.Model(&from).Association("Posts").Clear(); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err = tx.Model(&from).Association("Posts").Delete(fromRows); err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Model(to).Select("amount").Updates(map[string]interface{}{"amount": gorm.Expr("amount + ?", len(toRows))}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(toRows) > 0 {
		err = tx.Model(to).Association("Posts").Append(toRows)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return
}
