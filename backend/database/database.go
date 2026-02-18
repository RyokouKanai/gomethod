package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes the database connection.
func Connect() {
	host := getEnv("GMETHOD_DB_HOST", "db")
	user := getEnv("GMETHOD_DB_USERNAME", "root")
	pass := getEnv("GMETHOD_DB_PASSWORD", "password")
	dbName := getEnv("GMETHOD_DB_NAME", "gmethod_development")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Asia%%2FTokyo", user, pass, host, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
