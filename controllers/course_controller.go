package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Course struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

var database *gorm.DB
var databaseError error

func InitializeDatabase() {
	databaseHost := os.Getenv("DB_HOST")
	databasePort := os.Getenv("DB_PORT")
	databaseUser := os.Getenv("DB_USER")
	databaseName := os.Getenv("DB_NAME")
	databasePassword := os.Getenv("DB_PASSWORD")

	databaseURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, databaseUser, databaseName, databasePassword)
	fmt.Println("Connecting to database...")

	database, databaseError = gorm.Open("postgres", databaseURI)
	if databaseError != nil {
		panic("Failed to connect to the database")
	}

	database.AutoMigrate(&Course{})
}

func CreateCourseHandler(c *gin.Context) {
	var newCourse Course
	c.BindJSON(&newCourse)

	database.Create(&newCourse)
	c.IndentedJSON(http.StatusCreated, newCourse)
}

func GetAllCoursesHandler(c *gin.Context) {
	var courses []Course
	database.Find(&courses)
	c.IndentedJSON(http.StatusOK, courses)
}

func GetSingleCourseHandler(c *gin.Context) {
	courseID := c.Params.ByName("id")
	var course Course
	if err := database.Where("id = ?", courseID).First(&course).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, course)
}

func UpdateCourseHandler(c *gin.Context) {
	courseID := c.Params.ByName("id")
	var existingCourse Course
	if err := database.Where("id = ?", courseID).First(&existingCourse).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}

	var updatedCourseData Course
	c.BindJSON(&updatedCourseData)

	database.Model(&existingCourse).Updates(updatedCourseData)
	c.IndentedJSON(http.StatusOK, existingCourse)
}

func DeleteCourseHandler(c *gin.Context) {
	courseID := c.Params.ByName("id")
	var courseToDelete Course
	if err := database.Where("id = ?", courseID).First(&courseToDelete).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}

	database.Delete(&courseToDelete)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "course deleted"})
}

func SetupRouterEndpoints(router *gin.Engine) {
	router.GET("/courses", GetAllCoursesHandler)
	router.GET("/courses/:id", GetSingleCourseHandler)
	router.POST("/courses", CreateCourseHandler)
	router.PUT("/courses/:id", UpdateCourseHandler)
	router.DELETE("/courses/:id", DeleteCourseHandler)
}

func main() {
	InitializeDatabase()
	router := gin.Default()

	SetupRouterEndpoints(router)

	router.Run()
}