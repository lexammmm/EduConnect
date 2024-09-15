package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Title       string
	Description string
}

var DB *gorm.DB

func init() {
	loadEnv()
	setupDatabase()
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func setupDatabase() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	database.AutoMigrate(&Course{})
	DB = database
}

func main() {
	var course Course
	if err := DB.First(&course, 1).Error; err != nil {
		log.Printf("Error retrieving course: %v", err)
		return
	}

	fmt.Println("Course Title:", course.Title)
}