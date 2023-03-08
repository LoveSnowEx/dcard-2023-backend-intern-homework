package db

import (
	"fmt"

	"gorm.io/driver/postgres"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/config"
	"gorm.io/gorm"
)

var (
	db *DB
)

type DB struct {
	DB *gorm.DB
}

func Connect() (*DB, error) {
	if db != nil {
		return db, nil
	}

	// Get config
	host := config.DBHost
	port := config.DBPort
	user := config.DBUser
	password := config.DBPassword
	dbname := config.DBName
	timezone := config.TimeZone

	// Build data source name
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		host,
		user,
		password,
		dbname,
		port,
		timezone,
	)

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Migrate the schema
	db.AutoMigrate()

	return &DB{DB: db}, nil
}

func Close() error {
	if db == nil {
		return nil
	}

	dbConn, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	if err := dbConn.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
