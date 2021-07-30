package dao

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Database *gorm.DB

func Connect(env map[string]string) {
	databaseURL := env["DATABASE_URL"]
	var err error
	Database, err := gorm.Open("mysql", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	Database.LogMode(true)
	Database.AutoMigrate()
}
