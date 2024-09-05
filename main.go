package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sport_finder/api/handlers"
	"sport_finder/storage/postgres"
	"sport_finder/storage/redis"
)

const PORT = ":8080"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	db := postgres.NewPostgres()
	cache := redis.NewRedis()

	http.HandleFunc("POST /register", handlers.NewRegHandler(db))

	http.HandleFunc("POST /login", handlers.NewAuthHandler(db))

	http.HandleFunc("GET /objects", handlers.NewMapHandler(cache))

	log.Println("Starting server on ", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatalln("Error starting server", err)
	}
}
