package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() {
	var (
		dbUsername = os.Getenv("POSTGRES_USERNAME")
		dbPass     = os.Getenv("POSTGRES_PASSWORD")
		dbHost     = os.Getenv("POSTGRES_HOST")
		dbName     = os.Getenv("POSTGRES_DBNAME")
		dbPortStr  = os.Getenv("POSTGRES_PORT")
	)

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		panic(fmt.Sprintf("Failed to convert POSTGRES_PORT to an integer: %v", err))
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", dbHost, dbUsername, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to open DB: %v", err))
	}

	if gin.Mode() == gin.ReleaseMode {
		db.Logger.LogMode(0)
	}

	DB = db
	rawDB := RawDB()

	rawDB.SetMaxIdleConns(20)
	rawDB.SetMaxOpenConns(100)

	err = Migrate()
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate DB: %v", err))
	}
}

// RawDB returns the raw SQL database instance.
func RawDB() *sql.DB {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}

	return db
}
