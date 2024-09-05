package map_service

import (
	"encoding/json"
	"log"
	"sport_finder/services/map_service/osm"
)

type OSMResponse struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Type string  `json:"type"`
	ID   int64   `json:"id"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Tags Tags    `json:"tags"`
}

type Tags struct {
	Sport string `json:"sport"`
}

type Params struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Radius  int     `json:"radius"`
	Leisure string  `json:"leisure"`
}

// 10000,59.9343,30.3351
func GetAllObjects(params Params) ([]Element, error) {
	if params.Leisure == "" {
		params.Leisure = "pitch"
	}
	body, err := osm.GetObjectsFromOSM(params.Lat, params.Lon, params.Radius, params.Leisure)
	if err != nil {
		log.Printf("Failed to get objects from OSM: %v", err)
		return nil, err
	}
	var osmResponse OSMResponse
	err = json.Unmarshal(body, &osmResponse)
	if err != nil {
		log.Printf("Failed to unmarshal OSM response: %v", err)
		return nil, err
	}
	return osmResponse.Elements, nil
}
