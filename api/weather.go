package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type WeatherResponse struct {
	Wind_speed          int `json:"wind_speed"`
	Temperature_degrees int `json:"temperature_degrees"`
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func Get(query string) (*http.Response, error) {
	weatherstack_key := os.Getenv("WEATHERSTACK_KEY")
	request, err := http.NewRequest(http.MethodGet, "http://api.weatherstack.com/current?access_key="+weatherstack_key+"&query="+query, nil)
	if err != nil {
		return nil, err
	}
	return Client.Do(request)
}

func GetWeather(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		query = "Sydney"
	}
	log.Println("Looking up:", query)

	resp, err := Get(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		print(err)
	}

	wind_speed := data["current"].(map[string]interface{})["wind_speed"].(float64)
	temperature_degrees := data["current"].(map[string]interface{})["temperature"].(float64)
	weather_response := WeatherResponse{
		Wind_speed:          int(wind_speed),
		Temperature_degrees: int(temperature_degrees),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather_response)
}
