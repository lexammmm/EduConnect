package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type Course struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

var db *gorm.DB
var err error

type courseCache struct {
	courses    []Course
	lastUpdate time.Time
	mu         sync.RWMutex
}

var cCache courseCache

func initDB() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		panic(err)
	} else {
		fmt.Println("Database connected")
	}

	if err := db.AutoMigrate(&Course{}).Error; err != nil {
		fmt.Printf("Failed to migrate database: %v\n", err)
	}
}

func createCourse(c *gin.Context) {
	var course Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&course).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating course."})
		return
	}

	cCache.mu.Lock()
	cCache.courses = nil
	cCache.lastUpdate = time.Time{}
	cCache.mu.Unlock()

	c.JSON(http.StatusCreated, course)
}

func getCourses(c *gin.Context) {
	cCache.mu.RLock()
	if time.Since(cCache.lastUpdate) < 5*time.Minute && cCache.courses != nil {
		c.JSON(http.StatusOK, cCache.courses)
		cCache.mu.RUnlock()
		return
	}
	cCache.mu.RUnlock()

	var courses []Course
	if err := db.Find(&courses).Error; err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses."})
		return
	}

	cCache.mu.Lock()
	cCache.courses = courses
	cCache.lastUpdate = time.Now()
	cCache.mu.Unlock()

	c.JSON(http.StatusOK, courses)
}

func main() {
	initDB()

	r := gin.Default()

	r.GET("/courses", getCourses)
	r.POST("/course", createCourse)

	r.Run(":8080")
}