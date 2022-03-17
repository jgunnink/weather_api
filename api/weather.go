package api

import (
	"encoding/json"
	"errors"
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
	defer r.Body.Close()
	var wsr WeatherStackResponse
	var owmr OpenWeatherMapResponse

	query := r.URL.Query().Get("query")
	if query == "" {
		query = "Sydney"
	}

	clientResponse, err := callWeatherStack(query, wsr)
	if err != nil {
		// If we get a failure from WeatherStack, try OpenWeatherMap
		clientResponse, err = callOpenWeatherMap(query, owmr)
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientResponse)
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

	// Use the observation time since it will be there if there's a valid response. We can't check for zero value
	// temperature since it's possible that the temperature could be zero.
	if len(data.Current.ObservationTime) != 0 {
		return WeatherResponse{
			Wind_speed:          int(data.Current.WindSpeed),
			Temperature_degrees: int(data.Current.Temperature),
		}, nil
	} else {
		// Handle the case where WeatherStack returns a 200, but there is actually an error.
		var customResponse WeatherStackCustomResponse
		if err := json.NewDecoder(resp.Body).Decode(&customResponse); err != nil {
			return WeatherResponse{}, err
		}
		return WeatherResponse{}, errors.New(customResponse.Error.Info)
	}
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
