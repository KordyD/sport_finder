package map_service

import (
	"encoding/json"
	"log"
	"sport_finder/internal/services/map_service/osm"
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
type MapService struct{}

// 10000,59.9343,30.3351

func (mS *MapService) GetAllObjects(params Params) []Element {
	if params.Leisure == "" {
		params.Leisure = "pitch"
	}
	body := osm.GetObjectsFromOSM(params.Lat, params.Lon, params.Radius, params.Leisure)
	var osmResponse OSMResponse
	err := json.Unmarshal(body, &osmResponse)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	return osmResponse.Elements
}
