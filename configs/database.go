package configs

import (
	"easy-attend-service/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Database configuration from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	username := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	database := getEnv("DB_DATABASE", "easy-attend-serviceV2")

	// Create DSN (Data Source Name) for PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		host, username, password, database, port)

	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully!")

	// Auto migrate models
	AutoMigrate()
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.Gender{},
		&models.Prefix{},
		&models.School{},
		&models.Teacher{},
		&models.Student{},
		&models.Classroom{},
		&models.ClassroomMember{},
		&models.Attendance{},
		&models.Log{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed!")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
