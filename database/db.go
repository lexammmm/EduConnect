package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}

var DB *gorm.DB

func init() {
	loadEnv()
	setupDatabase()
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func setupDatabase() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	database.AutoMigrate(&Course{})
	DB = database
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/courses", getCourses).Methods("GET")
	api.HandleFunc("/course/{id}", getCourse).Methods("GET")
	api.HandleFunc("/course", addCourse).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getCourses(w http.ResponseWriter, r *http.Request) {
	var courses []Course
	DB.Find(&courses)
	json.NewEncoder(w).Encode(courses)
}

func getCourse(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	var course Course
	if err := DB.First(&course, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(course)
}

func addCourse(w http.ResponseWriter, r *http.Request) {
	var course Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	DB.Create(&course)
	json.NewEncoder(w).Encode(course)
}