package main

import (
	"github.com/gin-gonic/gin"
	"yourapp/controllers"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	router := gin.Default()
	router.POST("/courses", controllers.CreateCourse)
	router.GET("/courses/:id", controllers.GetCourse)
	router.PUT("/courses/:id", controllers.UpdateCourse)
	router.DELETE("/courses/:id", controllers.DeleteCourse)
	router.GET("/courses", controllers.GetAllCourses)
	router.Run(":" + getPort())
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}