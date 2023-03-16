package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/config"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/page"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/pagelist"
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable extension for uuid
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Migrate the schema
	db.AutoMigrate(page.Page{}, pagelist.PageNode{}, pagelist.PageList{})

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
