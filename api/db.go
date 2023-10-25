package api

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func InitDatabase() (*gorm.DB, error) {
	var DB *gorm.DB
	var err error
	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost != "" {
		DB, err = gorm.Open(postgres.Open("host=" + pgHost +
			" user=" + GetEnv("POSTGRES_USER", "postgres") +
			" password=" + GetEnv("POSTGRES_PASSWORD", "postgres") +
			" dbname=" + GetEnv("POSTGRES_DBNAME", "postgres") +
			" port=5432 sslmode=disable"))
		if err != nil {
			//log.Fatal("failed to initialize postgresql database")
			return nil, err
		}
	} else {
		log.Info("postgresql db not set, use sqlite file db")
		DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			//log.Fatal("failed to connect to sqlite3 database")
			return nil, err
		}
	}

	// create table and load data
	err = DB.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatal(err)
	}
	var cnt int64
	DB.Model(&Todo{}).Count(&cnt)
	if cnt == 0 {
		log.Info("Todos table empty, load initial data")
		DB.Create(&Todo{Title: "Test todo 1", Completed: false})
		DB.Create(&Todo{Title: "Test todo 2", Completed: false})
	}

	// sanity check
	var todo Todo
	DB.First(&todo)
	log.Infof("first todo in db: %v \n", todo)

	return DB, nil
}

type Todo struct {
	gorm.Model        // id, timestamping, and soft delete!
	Title      string `json:"title"`
	Completed  bool   `json:"completed"`
}
