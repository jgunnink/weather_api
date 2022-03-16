package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type WeatherResponse struct {
	Wind_speed          int `json:"wind_speed"`
	Temperature_degrees int `json:"temperature_degrees"`
}

func GetWeather(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	var weather_response WeatherResponse

	query := r.URL.Query().Get("query")
	if query == "" {
		query = "Sydney"
	}

	weather_response, err := tryWeatherStack(query, data)
	if err != nil {
		// If we get a failure from WeatherStack, try OpenWeatherMap
		weather_response, err = tryOpenWeatherMap(query, data)
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather_response)
}

func tryWeatherStack(query string, data map[string]interface{}) (WeatherResponse, error) {
	resp, err := GetFromWeatherStack(query)

	if err != nil {
		return WeatherResponse{}, err
	}
	if err := ValidateServiceResponse(resp); err != nil {
		return WeatherResponse{}, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherResponse{}, err
	}

	return setWeatherStack(resp, data)
}

func tryOpenWeatherMap(query string, data map[string]interface{}) (WeatherResponse, error) {
	resp, err := GetFromOpenWeatherMap(query)
	if err != nil {
		return WeatherResponse{}, err
	}
	err = ValidateServiceResponse(resp)
	if err != nil {
		return WeatherResponse{}, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherResponse{}, err
	}
	return setOpenWeather(resp, data)
}

func setWeatherStack(resp *http.Response, data map[string]interface{}) (WeatherResponse, error) {
	wind_speed := data["current"].(map[string]interface{})["wind_speed"].(float64)
	temperature_degrees := data["current"].(map[string]interface{})["temperature"].(float64)

	return WeatherResponse{
		Wind_speed:          int(wind_speed),
		Temperature_degrees: int(temperature_degrees),
	}, nil
}

func setOpenWeather(resp *http.Response, data map[string]interface{}) (WeatherResponse, error) {
	wind_speed := data["wind"].(map[string]interface{})["speed"].(float64)
	temperature_degrees := data["main"].(map[string]interface{})["temp"].(float64)

	return WeatherResponse{
		Wind_speed:          int(wind_speed),
		Temperature_degrees: int(temperature_degrees),
	}, nil
}
