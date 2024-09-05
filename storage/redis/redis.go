package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"sport_finder/services/map_service"
	"strconv"
	"time"
)

type Cache struct {
	RedisClient *redis.Client
}

func NewRedis() *Cache {
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
	return &Cache{RedisClient: rdb}
}

func (c *Cache) AddObjectsToCache(objects []map_service.Element) error {
	for _, object := range objects {
		bytes, err := json.Marshal(object)
		if err != nil {
			log.Println("Failed to marshal object: ", err)
			return err
		}
		status := c.RedisClient.Set(context.Background(), strconv.FormatInt(object.ID, 10), string(bytes), time.Second*15)
		if status.Err() != nil {
			log.Println("Failed to add object to cache: ", status.Err())
			return status.Err()
		}
		log.Println("Object added to cache")
	}
	return nil
}
