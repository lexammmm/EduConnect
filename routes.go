package main

import (
	"github.com/gin-gonic/gin"
	"yourapp/controllers"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// Optimized environment variables loading.
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load environment variables")
	}

	// Instance with default middleware.
	router := gin.Default()

	// Register endpoints
	registerRoutes(router)

	// Extracted Port configuration to reduce repetitive access to os.Getenv
	port := getPort()

	router.Run(":" + port)
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// Encapsulate route setup for potential future reuse and cleaner main function
func registerRoutes(router *gin.Engine) {
	router.POST("/courses", controllers.CreateCourse)
	router.GET("/courses/:id", controllers.GetCourse)
	router.PUT("/courses/:id", controllers.UpdateCourse)
	router.DELETE("/courses/:id", controllers.DeleteCourse)
	router.GET("/courses", controllers.GetAllCourses)
}