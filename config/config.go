package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var Default = Config{
	TimeZone:   "Asia/Taipei",
	DBHost:     "localhost",
	DBPort:     "5432",
	DBUser:     "postgres",
	DBPassword: "password",
	DBName:     "postgres",
	FiberPort:  "3000",
	GrpcPort:   "50051",
	GrpcuiPort: "8080",
}

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

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return filepath.Dir(filename)
}

func init() {
	cwd := getCurrentPath()
	rootwd := filepath.Join(cwd, "..")

	viper.AddConfigPath(rootwd)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
}

func Get() Config {
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("failed to read config: %v", err)
		return Default
	}

	return Config{
		TimeZone:   viper.GetString("TZ"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		FiberPort:  viper.GetString("FIBER_PORT"),
		GrpcPort:   viper.GetString("GRPC_PORT"),
		GrpcuiPort: viper.GetString("GRPCUI_PORT"),
	}
}
