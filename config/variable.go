package config

import "gorm.io/gorm"

var (
	appPort string
	appHost string
	// pythonNodePort string

	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string

	dbPsql *gorm.DB
)
