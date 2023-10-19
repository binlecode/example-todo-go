package main

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is exported global variable to hold the database connection pool.
var DB *gorm.DB

func initDatabase() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to sqlite3 database")
	}

	// create table and load data
	err = DB.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatal(err)
	}
	DB.Create(&Todo{Title: "Test todo 1", Completed: false})
	DB.Create(&Todo{Title: "Test todo 2", Completed: false})

	// sanity check
	var todo Todo
	DB.First(&todo, 1)
	log.Infof("first loaded todo: %v \n", todo)

	return nil
}

type Todo struct {
	gorm.Model        // id, timestamping, and soft delete!
	Title      string `json:"title"`
	Completed  bool   `json:"completed"`
}
