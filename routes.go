package main

import (
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "os"
    "yourapp/controllers"

    "github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Failed to load environment variables: %v", err)
    }

    router := gin.Default()

    registerRoutes(router)

    port := getPort()

    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

func getPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    return port
}

func registerRoutes(router *gin.Engine) {
    router.POST("/courses", controllers.CreateCourse)
    router.GET("/courses/:id", controllers.GetCourse)
    router.PUT("/courses/:id", controllers.UpdateCourse)
    router.DELETE("/courses/:id", controllers.DeleteCourse)
    router.GET("/courses", controllers.GetAllCourses)
}