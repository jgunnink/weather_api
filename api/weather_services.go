package api

import (
	"log"
	"net/http"
	"os"
)

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

type WeatherResponseErr string

func (e WeatherResponseErr) Error() string {
	return string(e)
}

const (
	ErrNotFound = WeatherResponseErr("The query for that location returned an empty result")
	ErrUpstream = WeatherResponseErr("The upstream service returned an error")
	ErrApiKey   = WeatherResponseErr("The weather service is busy right now try again soon")
)

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

func ValidateServiceResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return ErrNotFound
		case http.StatusTooManyRequests:
			return ErrApiKey
		default:
			return ErrUpstream
		}
	}
	return nil
}
