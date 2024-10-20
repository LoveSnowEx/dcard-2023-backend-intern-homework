package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"

	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/config"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/page"
	"github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/pagelist"
	"gorm.io/gorm"
)

var (
	db         *DB
	gormConfig = gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Warn,
				Colorful:      true,
			},
		),
	}
)

type DB struct {
	DB *gorm.DB
}

func Connect() (*DB, error) {
	if db != nil {
		return db, nil
	}

	var dialector gorm.Dialector

	conf := config.Get()
	switch conf.DBDriver {
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(conf.DBName)
	case "postgres":
		dialector = postgres.Open(fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			conf.DBHost,
			conf.DBUser,
			conf.DBPassword,
			conf.DBName,
			conf.DBPort,
			conf.DBSSLMode,
			conf.TimeZone,
		))
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", conf.DBDriver)
	}

	conn, err := gorm.Open(dialector, &gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable extension for uuid
	conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Migrate the schema
	if err := conn.AutoMigrate(page.Page{}, pagelist.PageNode{}, pagelist.PageList{}); err != nil {
		log.Printf("failed to migrate schema: %v", err)
	}

	db = &DB{DB: conn}

	return db, nil
}

func Close() error {
	if db == nil {
		return nil
	}

	conn, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	err = conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	db = nil

	return nil
}
