package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Course{})
}

type Course struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"size:1024;not null"`
	Instructor  string `gorm:"size:255;not null"`
}

func main() {
}