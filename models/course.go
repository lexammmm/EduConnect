package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
)

var databaseConnection *gorm.DB

func init() {
    databaseURL := os.Getenv("DATABASE_URL")
    var err error
    databaseConnection, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
    if err != nil {
        panic("failed to connect to the database")
    }

    err = databaseConnection.AutoMigrate(&Course{})
    if err != nil {
        panic("auto migration failed")
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

func closeDatabaseConnection() {
    sqlDB, err := databaseConnection.DB()
    if err != nil {
        panic("failed to retrieve the SQL database object from GORM")
    }
    if err := sqlDB.Close(); err != nil {
        panic("failed to close the SQL database connection")
    }
}