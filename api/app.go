package api

import "gorm.io/gorm"

// App is application object to hold the dependencies

type App struct {
	DB *gorm.DB
}
