package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./local.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	return db, err
}