package controllers

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"net/http"
	"sport_finder/internal/services/map_service"
	"sport_finder/internal/storage"
)

type MapObjectsProvider interface {
	GetAllObjects(params map_service.Params) []map_service.Element
}

func MapController(w http.ResponseWriter, r *http.Request, mapObjectsProvider MapObjectsProvider, rdb *redis.Client) {
	var mapParams map_service.Params
	err := json.NewDecoder(r.Body).Decode(&mapParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	objects := mapObjectsProvider.GetAllObjects(mapParams)
	storage.AddArrayOfObjects(objects, storage.Cache{RedisClient: rdb})
	err = json.NewEncoder(w).Encode(objects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
