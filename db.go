package main

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDatabase() error {
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to sqlite3 database")
	}

	// create table and load data
	db.AutoMigrate(&Todo{})
	db.Create(&Todo{Title: "Test todo 1", Completed: false})
	db.Create(&Todo{Title: "Test todo 2", Completed: false})

	// sanity check
	var todo Todo
	db.First(&todo, 1)
	log.Infof("first loaded todo: %v \n", todo)

	return nil
}
