package database

import (
	"fmt"
	"log"
	"os"
	"product-service/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDB initializes the database connection and assigns it to DB
func ConnectDB() {
	err := godotenv.Load(".env") // Ensure the correct path
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Connected to database successfully")

	// Assign to global DB variable
	DB = db

	// Migrate models
	err = DB.AutoMigrate(
		&models.Product{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	} else {
		fmt.Println("✅ Migrations completed successfully")
	}
}
