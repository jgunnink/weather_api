package api

import (
	"encoding/json"
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
	var clientResponse WeatherResponse
	var err error
	var providers = [...]string{"WeatherStack", "OpenWeatherMap"}

	query := r.URL.Query().Get("query")
	if query == "" {
		query = "Sydney"
	}

	validResponse := false
	for _, provider := range providers {
		if validResponse {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(clientResponse)
			return
		}
		switch provider {
		case "WeatherStack":
			clientResponse, err = callWeatherStack(query, wsr)
		case "OpenWeatherMap":
			clientResponse, err = callOpenWeatherMap(query, owmr)
		default:
			http.Error(w, "Error: invalid upstream service provider.", http.StatusInternalServerError)
		}

		if err == nil {
			validResponse = true
		} else {
			http.Error(w,
				"Error: invalid response from our upstream weather providers. Please a different query or try again later.",
				http.StatusServiceUnavailable)
		}
	}
}

func callWeatherStack(query string, data WeatherStackResponse) (WeatherResponse, error) {
	resp, err := GetFromWeatherStack(query)
	if err != nil {
		return WeatherResponse{}, err
	}
	if err := ValidateUpstreamResponse(resp); err != nil {
		return WeatherResponse{}, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherResponse{}, err
	}

	return WeatherResponse{
		Wind_speed:          data.Current.WindSpeed,
		Temperature_degrees: data.Current.Temperature,
	}, nil
}

func callOpenWeatherMap(query string, data OpenWeatherMapResponse) (WeatherResponse, error) {
	resp, err := GetFromOpenWeatherMap(query)
	if err != nil {
		return WeatherResponse{}, err
	}
	if err := ValidateUpstreamResponse(resp); err != nil {
		return WeatherResponse{}, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return WeatherResponse{}, err
	}

	return WeatherResponse{
		Wind_speed:          int(data.Wind.Speed),
		Temperature_degrees: int(data.Main.Temp),
	}, nil
}
