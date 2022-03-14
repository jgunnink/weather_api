package main

import (
	"net/http"

	"github.com/jgunnink/weather_api/api"
)

func main() {
	http.HandleFunc("/v1/weather", func(w http.ResponseWriter, r *http.Request) {
		api.GetWeather(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
