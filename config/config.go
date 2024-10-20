package config

import (
	"log"

	"github.com/spf13/viper"
)

var config Config

type Config struct {
	TimeZone   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	FiberPort  string
	GrpcPort   string
	GrpcuiPort string
}

func init() {
	viper.SetConfigFile(".env")

	viper.SetDefault("TZ", "Asia/Taipei")
	viper.SetDefault("DB_HOST", "postgres")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "postgres")
	viper.SetDefault("FIBER_PORT", "3000")
	viper.SetDefault("GRPC_PORT", "50051")
	viper.SetDefault("GRPCUI_PORT", "8080")
}

func Get() Config {
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("failed to read config: %v", err)
		return config
	}

	config.TimeZone = viper.GetString("TZ")
	config.DBUser = viper.GetString("DB_USER")
	config.DBPassword = viper.GetString("DB_PASSWORD")
	config.DBName = viper.GetString("DB_NAME")
	config.DBHost = viper.GetString("DB_HOST")
	config.DBPort = viper.GetString("DB_PORT")
	config.FiberPort = viper.GetString("FIBER_PORT")
	config.GrpcPort = viper.GetString("GRPC_PORT")
	config.GrpcuiPort = viper.GetString("GRPCUI_PORT")

	return config
}
