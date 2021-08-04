package model

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(env map[string]string) {
	dsn := env["DATABASE_URL"]
	var err error
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatal(err)
	}
	// db.Debug().Logger
	db.AutoMigrate(&User{})
}

func applyQueryOptions(options map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if options["preload"] != nil {
			if preload, ok := options["preload"].([]string); ok {
				for _, col := range preload {
					db = db.Preload(col)
				}
			}
			if preload, ok := options["preload"].(map[string]interface{}); ok {
				for key, val := range preload {
					if val == nil {
						db = db.Preload(key)
					} else {
						db = db.Preload(key, val)
					}
				}
			}
		}
		if options["select"] != nil {
			if selected, ok := options["select"].([]string); ok {
				db = db.Select(selected)
			}
		}
		if options["where"] != nil {
			if wheres, ok := options["where"].([][]interface{}); ok {
				for _, where := range wheres {
					db = db.Where(where[0], where[1:]...)
				}
			}
			if wheres, ok := options["where"].(map[string]interface{}); ok {
				db = db.Where(wheres)
			}
		}
		if options["join"] != nil {
			db = db.Joins(options["join"].(string))
		}
		if options["order"] != nil {
			if orders, ok := options["order"].([]string); ok {
				for _, order := range orders {
					db = db.Order(order)
				}
			} else {
				db = db.Order(options["order"])
			}
		}
		if options["offset"] != nil {
			db = db.Offset(options["offset"].(int))
		}
		if options["limit"] != nil {
			db = db.Limit(options["limit"].(int))
		}
		return db
	}
}
