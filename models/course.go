package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
)

var db *gorm.DB

func init() {
    dbURL := os.Getenv("DATABASE_URL")
    var err error
    db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    err = db.AutoMigrate(&Course{})
    if err != nil {
        panic("failed to auto migrate")
    }
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

func closeDB() {
    sqlDB, err := db.DB()
    if err != nil {
        panic("failed to get database object from GORM connection")
    }
    if err := sqlDB.Close(); err != nil {
        panic("failed to close database connection")
    }
}