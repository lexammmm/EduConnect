package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"fmt"
	"net/http"
)

type Course struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

var db *gorm.DB

func initDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	fmt.Println("Connecting to database...")

	var err error
	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		panic("Failed to connect to the database")
	}

	db.AutoMigrate(&Course{})
}

func createCourse(c *gin.Context) {
	var course Course
	if err := c.BindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	db.Create(&course)
	c.JSON(http.StatusCreated, course)
}

func getAllCourses(c *gin.Context) {
	var courses []Course
	db.Find(&courses)
	c.JSON(http.StatusOK, courses)
}

func getCourse(c *gin.Context) {
	var course Course
	if err := db.First(&course, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

func updateCourse(c *gin.Context) {
	var course Course
	if err := db.First(&course, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	var update Course
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	db.Model(&course).Updates(update)
	c.JSON(http.StatusOK, course)
}

func deleteCourse(c *gin.Context) {
	var course Course
	if err := db.Where("id = ?", c.Param("id")).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	db.Delete(&course)
	c.JSON(http.StatusOK, gin.H{"message": "course deleted"})
}

func setupRoutes(router *gin.Engine) {
	router.GET("/courses", getAllCourses)
	router.GET("/courses/:id", getCourse)
	router.POST("/courses", createCourse)
	router.PUT("/courses/:id", updateCourse)
	router.DELETE("/courses/:id", deleteCourse)
}

func main() {
	initDB()
	r := gin.Default()

	setupRoutes(r)

	r.Run()
}