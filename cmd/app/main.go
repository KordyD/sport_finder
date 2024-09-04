package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sport_finder/internal/routes"
	"sport_finder/internal/storage/redis"
)

const PORT = ":8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connection := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", dbUsername, dbPassword, dbName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatalln("Error in connection", err)
	}
	rdb := redis.Connect()
	routes.RegisterRoutes(db, rdb)

	log.Println("Starting server on ", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatalln("Error starting server", err)
	}
}
