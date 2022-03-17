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
	var wsr WeatherStackResponse
	var owmr OpenWeatherMapResponse
	var client_response WeatherResponse

	query := r.URL.Query().Get("query")
	if query == "" {
		query = "Sydney"
	}

	client_response, err := tryWeatherStack(query, wsr)
	if err != nil {
		// If we get a failure from WeatherStack, try OpenWeatherMap
		client_response, err = tryOpenWeatherMap(query, owmr)
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client_response)
}

func tryWeatherStack(query string, data WeatherStackResponse) (WeatherResponse, error) {
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

	return setWeatherStack(resp, data), nil
}

func tryOpenWeatherMap(query string, data OpenWeatherMapResponse) (WeatherResponse, error) {
	resp, err := GetFromOpenWeatherMap(query)
	if err != nil {
		return WeatherResponse{}, err
	}
	if err := ValidateServiceResponse(resp); err != nil {
		return WeatherResponse{}, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherResponse{}, err
	}
	return setOpenWeather(resp, data), nil
}

func setWeatherStack(resp *http.Response, data WeatherStackResponse) WeatherResponse {
	wind_speed := data.Current.WindSpeed
	temperature_degrees := data.Current.Temperature

	return WeatherResponse{
		Wind_speed:          int(wind_speed),
		Temperature_degrees: int(temperature_degrees),
	}
}

func setOpenWeather(resp *http.Response, data OpenWeatherMapResponse) WeatherResponse {
	wind_speed := data.Wind.Speed
	temperature_degrees := data.Wind.Speed

	return WeatherResponse{
		Wind_speed:          int(wind_speed),
		Temperature_degrees: int(temperature_degrees),
	}
}
