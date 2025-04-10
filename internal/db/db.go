package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() {

	// try to connect to neon db
	if connectToRemoteDB() {
		log.Println("Connected to intercordDB successfully")
	} else {
		// Fall back to local DB if remote connection fails
		log.Println("Remote database connection failed, attempting to connect to local database")
		connectToLocalDB()
	}

	rawDB := RawDB()

	rawDB.SetMaxIdleConns(20)
	rawDB.SetMaxOpenConns(100)

	err := Migrate()
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate DB: %v", err))
	}
}

func connectToRemoteDB() bool {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL variable not found in .env")
		return false
	}

	// open connection
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to intercordDB: %v", err)
		return false
	}

	// test & ping the connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed  to get connection: %v", err)
		return false
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Failed to ping intercordDB: %v", err)
		return false
	}

	DB = db // set the global var

	// Configure logging level based on Gin mode
	if gin.Mode() == gin.ReleaseMode {
		db.Logger.LogMode(0)
	}

	return true
}

func connectToLocalDB() {
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
}

// RawDB returns the raw SQL database instance.
func RawDB() *sql.DB {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}

	return db
}
