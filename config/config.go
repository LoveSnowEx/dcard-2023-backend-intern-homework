package config

import (
	"github.com/spf13/viper"
)

var config Config

type Config struct {
	TimeZone   string
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	FiberPort  string
	GrpcPort   string
	GrpcuiPort string
}

func init() {
	viper.AutomaticEnv()
	viper.SetDefault("TZ", "Asia/Taipei")
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_HOST", "postgres")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "postgres")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("FIBER_PORT", "3000")
	viper.SetDefault("GRPC_PORT", "50051")
	viper.SetDefault("GRPCUI_PORT", "8080")

	config.TimeZone = viper.GetString("TZ")
	config.DBDriver = viper.GetString("DB_DRIVER")
	config.DBUser = viper.GetString("DB_USER")
	config.DBPassword = viper.GetString("DB_PASSWORD")
	config.DBName = viper.GetString("DB_NAME")
	config.DBHost = viper.GetString("DB_HOST")
	config.DBPort = viper.GetString("DB_PORT")
	config.DBSSLMode = viper.GetString("DB_SSLMODE")
	config.FiberPort = viper.GetString("FIBER_PORT")
	config.GrpcPort = viper.GetString("GRPC_PORT")
	config.GrpcuiPort = viper.GetString("GRPCUI_PORT")
}

func Get() Config {
	return config
}
