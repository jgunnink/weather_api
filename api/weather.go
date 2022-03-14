package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func GetWeather(w http.ResponseWriter, r *http.Request) {
	weatherstack_key := os.Getenv("WEATHERSTACK_KEY")
	resp, err := http.Get("http://api.weatherstack.com/current?access_key=" + weatherstack_key + "&query=Sydney")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	w.Write(bodyBytes)

}
