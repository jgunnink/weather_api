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

func GetFromWeatherStack(query string) (*http.Response, error) {
	log.Println("Looking up:", query, "on Weather Stack")
	weatherstack_key := os.Getenv("WEATHERSTACK_KEY")
	request, err := http.NewRequest(http.MethodGet, "http://api.weatherstack.com/current?access_key="+weatherstack_key+"&query="+query, nil)
	if err != nil {
		return nil, err
	}
	return Client.Do(request)
}

func GetFromOpenWeatherMap(query string) (*http.Response, error) {
	log.Println("Looking up:", query, "on Open Weather Map")
	openweathermap_key := os.Getenv("OPENWEATHERMAP_KEY")
	request, err := http.NewRequest(http.MethodGet, "http://api.openweathermap.org/data/2.5/weather?q="+query+"&units=metric&appid="+openweathermap_key, nil)
	if err != nil {
		return nil, err
	}
	return Client.Do(request)
}

func GetWeather(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	var weather_response WeatherResponse

	query := r.URL.Query().Get("query")
	if query == "" {
		query = "Sydney"
	}

	resp, err := GetFromWeatherStack(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		return
	}
	weather_response, err = setWeatherStack(resp, data)

	// If we get a failure from WeatherStack, try OpenWeatherMap
	if err != nil {
		resp, err := GetFromOpenWeatherMap(query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
			return
		}

		weather_response, err = setOpenWeather(resp, data)
		if err != nil {
			// If we get two invalid responses from our upstream, then we should respond to our user with a 502
			http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusBadGateway)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather_response)
}

func setWeatherStack(resp *http.Response, data map[string]interface{}) (WeatherResponse, error) {
	if resp.Body == nil {
		return WeatherResponse{}, fmt.Errorf("no response body from WeatherStack")
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&data)
	wind_speed := data["current"].(map[string]interface{})["wind_speed"].(float64)
	temperature_degrees := data["current"].(map[string]interface{})["temperature"].(float64)

	return WeatherResponse{
		Wind_speed:          int(wind_speed),
		Temperature_degrees: int(temperature_degrees),
	}, nil
}

func setOpenWeather(resp *http.Response, data map[string]interface{}) (WeatherResponse, error) {
	if resp.Body == nil {
		return WeatherResponse{}, fmt.Errorf("no response body from OpenWeather")
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&data)
	wind_speed := data["wind"].(map[string]interface{})["speed"].(float64)
	temperature_degrees := data["main"].(map[string]interface{})["temp"].(float64)

	return WeatherResponse{
		Wind_speed:          int(wind_speed),
		Temperature_degrees: int(temperature_degrees),
	}, nil
}
