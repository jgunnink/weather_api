package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	weatherstack_key := os.Getenv("WEATHERSTACK_KEY")
	http.HandleFunc("/v1/weather", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://api.weatherstack.com/current?access_key=" + weatherstack_key + "&query=Sydney")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
		}
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		w.Write(bodyBytes)
	})

	http.ListenAndServe(":8080", nil)
}
