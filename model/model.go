package model

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func InitDB(env map[string]string) {
	conn := env["DATABASE_URL"]
	var err error
	db, err = gorm.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	db.AutoMigrate(&User{})
}

// func applyQueryOptions()
