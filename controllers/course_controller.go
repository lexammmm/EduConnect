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

var db *gorm.DB
var err error

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dbUri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
	fmt.Println("Connecting to database...")

	db, err = gorm.Open("postgres", dbUri)
	if err != nil {
		panic("Failed to connect to the database")
	}

	db.AutoMigrate(&Course{})
}

func CreateCourse(c *gin.Context) {
	var course Course
	c.BindJSON(&course)

	db.Create(&course)
	c.IndentedJSON(http.StatusCreated, course)
}

func GetCourses(c *gin.Context) {
	var courses []Course
	db.Find(&courses)
	c.IndentedJSON(http.StatusOK, courses)
}

func GetCourse(c *gin.Context) {
	id := c.Params.ByName("id")
	var course Course
	if err := db.Where("id = ?", id).First(&course).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, course)
}

func UpdateCourse(c *gin.Context) {
	var course Course
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&course).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}

	var newCourse Course
	c.BindJSON(&newCourse)

	db.Model(&course).Updates(newCourse)
	c.IndentedJSON(http.StatusOK, course)
}

func DeleteCourse(c *gin.Context) {
	id := c.Params.ByName("id")
	var course Course
	if err := db.Where("id = ?", id).First(&course).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "course not found"})
		return
	}

	db.Delete(&course)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "course deleted"})
}

func SetupRoutes(r *gin.Engine) {
	r.GET("/courses", GetCourses)
	r.GET("/courses/:id", GetCourse)
	r.POST("/courses", CreateCourse)
	r.PUT("/courses/:id", UpdateCourse)
	r.DELETE("/courses/:id", DeleteCourse)
}

func main() {
	InitDB()
	r := gin.Default()

	SetupRoutes(r)

	r.Run()
}