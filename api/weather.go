package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
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
	resp, err := Get("Sydney")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	w.Write(bodyBytes)

}
