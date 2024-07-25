package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
	"sync"
	"time"
)

type Course struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

var db *gorm.DB
var err error

// Cache structure
type courseCache struct {
	courses []Course
	lastUpdate time.Time
	mu sync.RWMutex
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
		fmt.Println("Failed to connect to database")
		panic(err)
	} else {
		fmt.Println("Database connected")
	}

	db.AutoMigrate(&Course{})
}

func createCourse(c *gin.Context) {
	var course Course
	c.BindJSON(&course)
	db.Create(&course)
	// Invalidate the cache
	cCache.mu.Lock()
	cCache.courses = nil // Clear the cache to force a refresh on next read
	cCache.mu.Unlock()
	c.JSON(200, course)
}

func getCourses(c *gin.Context) {
	cCache.mu.RLock()
	if time.Since(cCache.lastUpdate) < 5*time.Minute && cCache.courses != nil {
		// Serve from cache
		c.JSON(200, cCache.courses)
		cCache.mu.RUnlock()
		return
	}
	cCache.mu.RUnlock()

	var courses []Course
	if err := db.Find(&courses).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		cCache.mu.Lock()
		cCache.courses = courses
		cCache.lastUpdate = time.Now()
		cCache.mu.Unlock()
		c.JSON(200, courses)
	}
}

func main() {
	initDB()

	r := gin.Default()

	r.GET("/courses", getCourses)
	r.POST("/course", createCourse)

	r.Run(":8080")
}