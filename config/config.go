package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

var (
	TimeZone   string = "Asia/Taipei"
	DBHost     string = "localhost"
	DBPort     int    = 5432
	DBUser     string = "postgres"
	DBPassword string = "password"
	DBName     string = "postgres"
	FiberPort  string = "3000"
	GrpcPort   string = "50051"
	GrpcuiPort string = "8080"
)

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)

	return filepath.Dir(filename)
}

func init() {
	cwd := getCurrentPath()
	rootwd := filepath.Join(cwd, "..")

	viper.SetConfigType("env")
	viper.AddConfigPath(rootwd)
	viper.SetConfigName(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("failed to read config: %v", err)
		return
	}

	TimeZone = viper.GetString("TZ")

	DBHost = viper.GetString("DB_HOST")
	DBPort = viper.GetInt("DB_PORT")
	DBUser = viper.GetString("DB_USER")
	DBPassword = viper.GetString("DB_PASSWORD")
	DBName = viper.GetString("DB_NAME")

	FiberPort = viper.GetString("FIBER_PORT")
	GrpcPort = viper.GetString("GRPC_PORT")
	GrpcuiPort = viper.GetString("GRPCUI_PORT")
}
