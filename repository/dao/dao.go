package dao

import (
	"app/util"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init(dsn string) {
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	// db.Debug().Logger
	db.AutoMigrate(&User{}, &Post{}, &Category{})
}

func Close() error {
	d, err := db.DB()
	if err != nil {
		return err
	}
	return d.Close()
}

type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt util.LocalTime `json:"createdAt"`
	UpdatedAt util.LocalTime `json:"updatedAt"`
	DeletedAt util.DeletedAt `gorm:"index" json:"deletedAt"`
}

func applyQueryOptions(options map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// tx := db.Session(&gorm.Session{})
		if options["preload"] != nil {
			if preload, ok := options["preload"].([]string); ok {
				for _, col := range preload {
					tx = tx.Preload(col)
				}
			}
			if preload, ok := options["preload"].(map[string]interface{}); ok {
				for key, val := range preload {
					if val == nil {
						tx = tx.Preload(key)
					} else {
						tx = tx.Preload(key, val)
					}
				}
			}
		}
		if options["select"] != nil {
			if selected, ok := options["select"].([]string); ok {
				tx = tx.Select(selected)
			}
		}
		if options["where"] != nil {
			switch options["where"].(type) {
			case []string, []uint:
				tx = tx.Where("id in (?)", options["where"])
			case [][]interface{}:
				for _, where := range options["where"].([][]interface{}) {
					tx = tx.Where(where[0], where[1:]...)
				}
			case map[string]interface{}:
				tx = tx.Where(options["where"])
			}
		}
		if options["join"] != nil {
			tx = tx.Joins(options["join"].(string))
		}
		if options["order"] != nil {
			if orders, ok := options["order"].([]string); ok {
				for _, order := range orders {
					tx = tx.Order(order)
				}
			} else {
				tx = tx.Order(options["order"])
			}
		}
		if options["offset"] != nil {
			tx = tx.Offset(options["offset"].(int))
		}
		if options["limit"] != nil {
			tx = tx.Limit(options["limit"].(int))
		}
		return tx
	}
}

func initData() {
	var category = Category{
		Name: "根分类", Description: "根分类", Lft: 1, Rgt: 2, Depth: 0,
	}
	_, err := category.Create(nil)
	if err != nil {
		log.Fatal(err)
	}
}
