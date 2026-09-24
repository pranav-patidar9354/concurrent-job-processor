package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pranav-patidar9354/concurrent-job-processor/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	// Load environment variables
	_ = godotenv.Load()

	// Prioritize DATABASE_URL (for Neon / Cloud PostgreSQL)
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}

		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}

		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "job_processor"
		}

		sslmode := os.Getenv("DB_SSLMODE")
		if sslmode == "" {
			sslmode = "disable"
		}

		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host,
			port,
			user,
			password,
			dbname,
			sslmode,
		)
	}

	// Connect GORM to PostgreSQL
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = database

	fmt.Println("PostgreSQL database connected successfully!")

	// Automatically create/update database tables
	err = DB.AutoMigrate(
		&models.Job{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	fmt.Println("Database migration completed successfully!")
}