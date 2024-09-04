package storage

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
	"sport_finder/internal/services/map_service"
	"strconv"
	"time"
)

type Cache struct {
	RedisClient *redis.Client
}

func AddArrayOfObjects(objects []map_service.Element, cacheDb Cache) {
	for _, object := range objects {
		bytes, err := json.Marshal(object)
		if err != nil {
			log.Println("Failed to marshal object: ", err)
		}
		status := cacheDb.RedisClient.Set(context.Background(), strconv.FormatInt(object.ID, 10), string(bytes), time.Second*15)
		if status.Err() != nil {
			log.Println("Failed to add object to cache: ", status.Err())
		}
		log.Println("Object added to cache")
	}
}
