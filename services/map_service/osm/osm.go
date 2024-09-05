package osm

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const apiUrl = "http://overpass-api.de/api/interpreter"

func GetObjectsFromOSM(lat float64, lon float64, radiusInMeters int, leisure string) ([]byte, error) {
	query := fmt.Sprintf(`[out:json];node["leisure"="%s"](around:%d,%f,%f);out body;`, leisure, radiusInMeters, lat, lon)
	base, err := url.Parse(apiUrl)
	if err != nil {
		log.Println("Failed to parse URL: ", err)
		return nil, err
	}
	params := url.Values{}
	params.Add("data", query)
	base.RawQuery = params.Encode()
	resp, err := http.Get(base.String())
	if err != nil {
		log.Printf("Failed to get data from Overpass API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}
	return body, nil
}
