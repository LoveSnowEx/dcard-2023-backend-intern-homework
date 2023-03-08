package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	TimeZone   string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	TimeZone = viper.GetString("TZ")

	DBHost = viper.GetString("DB_HOST")
	DBPort = viper.GetInt("DB_PORT")
	DBUser = viper.GetString("DB_USER")
	DBPassword = viper.GetString("DB_PASSWORD")
	DBName = viper.GetString("DB_NAME")
}
