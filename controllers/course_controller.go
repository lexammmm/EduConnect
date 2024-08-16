package main

import (
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

func initializeDatabase() {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbName := os.Getenv("DB_NAME")
    dbPassword := os.Getenv("DB_PASSWORD")

    databaseConnectionURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword)
    fmt.Println("Connecting to database...")

    var err error
    database, err = gorm.Open("postgres", databaseConnectionURI)
    if err != nil {
        panic("Failed to connect to the database")
    }

    database.AutoMigrate(&Course{})
}

func handleCreateCourse(c *gin.Context) {
    var newCourse Course
    if err := c.BindJSON(&newCourse); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    database.Create(&newCourse)
    c.JSON(http.StatusCreated, newCourse)
}

func handleGetAllCourses(c *gin.Context) {
    var courses []Course
    database.Find(&courses)
    c.JSON(http.StatusOK, courses)
}

func handleGetCourse(c *gin.Context) {
    var course Course
    if err := database.First(&course, "id = ?", c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
        return
    }
    c.JSON(http.StatusOK, course)
}

func handleUpdateCourse(c *gin.Context) {
    var course Course
    if err := database.First(&course, "id = ?", c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
        return
    }

    var updatedCourse Course
    if err := c.BindJSON(&updatedCourse); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
        return
    }

    database.Model(&course).Updates(updatedCourse)
    c.JSON(http.StatusOK, course)
}

func handleDeleteCourse(c *gin.Context) {
    var course Course
    if err := database.Where("id = ?", c.Param("id")).First(&course).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
        return
    }

    database.Delete(&course)
    c.JSON(http.StatusOK, gin.H{"message": "course deleted"})
}

func setupRouter(router *gin.Engine) {
    router.GET("/courses", handleGetAllCourses)
    router.GET("/courses/:id", handleGetCourse)
    router.POST("/courses", handleCreateCourse)
    router.PUT("/courses/:id", handleUpdateCourse)
    router.DELETE("/courses/:id", handleDeleteCourse)
}

func main() {
    initializeDatabase()
    router := gin.Default()

    setupRouter(router)

    router.Run()
}