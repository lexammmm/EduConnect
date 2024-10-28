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

var databaseInstance *gorm.DB

func init() {
    loadEnvironmentVariables()
    initializeDatabaseConnection()
}

func loadEnvironmentVariables() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file:", err)
    }
}

func initializeDatabaseConnection() {
    username := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASS")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")

    dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, username, dbName, password)

    database, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    database.AutoMigrate(&Course{})
    databaseInstance = database
}

func main() {
    router := mux.NewRouter()
    apiRouter := router.PathPrefix("/api/v1").Subrouter()

    apiRouter.HandleFunc("/courses", retrieveCoursesHandler).Methods("GET")
    apiRouter.HandleFunc("/course/{id}", retrieveCourseHandler).Methods("GET")
    apiRouter.HandleFunc("/course", createCourseHandler).Methods("POST")

    log.Fatal(http.ListenAndServe(":8080", router))
}

func retrieveCoursesHandler(responseWriter http.ResponseWriter, request *http.Request) {
    var courses []Course
    databaseInstance.Find(&courses)
    respondWithJSON(responseWriter, courses)
}

func retrieveCourseHandler(responseWriter http.ResponseWriter, request *http.Request) {
    parameters := mux.Vars(request)
    id, _ := strconv.Atoi(parameters["id"])
    var course Course
    if err := databaseInstance.First(&course, id).Error; err != nil {
        http.Error(responseWriter, "Course not found", http.StatusNotFound)
        return
    }
    respondWithJSON(responseWriter, course)
}

func createCourseHandler(responseWriter http.ResponseWriter, request *http.Request) {
    var course Course
    if err := json.NewDecoder(request.Body).Decode(&course); err != nil {
        http.Error(responseWriter, "Invalid request payload", http.StatusBadRequest)
        return
    }
    databaseInstance.Create(&course)
    respondWithJSON(responseWriter, course)
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
    json.NewEncoder(w).Encode(data)
}