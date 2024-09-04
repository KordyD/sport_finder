package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func Connect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "",
		DB:       0,
	})
	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		log.Fatalln("Failed to connect to Redis: ", status.Err())
	}
	log.Println("Connected to Redis")
	return rdb
}
