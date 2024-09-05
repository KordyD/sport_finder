package handlers

import (
	"encoding/json"
	"net/http"
	"sport_finder/services/map_service"
)

type MapObjectsProvider interface {
	AddObjectsToCache(objects []map_service.Element) error
}

func NewMapHandler(mapObjectsProvider MapObjectsProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mapParams map_service.Params
		err := json.NewDecoder(r.Body).Decode(&mapParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		objects, err := map_service.GetAllObjects(mapParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = mapObjectsProvider.AddObjectsToCache(objects)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(objects)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
