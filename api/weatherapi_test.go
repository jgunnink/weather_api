package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jgunnink/weather_api/utils/mocks"
)

func init() {
	Client = &mocks.MockClient{}
}

func TestGetWeather(t *testing.T) {
	Client = &mocks.MockClient{}
	standard_response := `{"request":{"type":"City","query":"Sydney, Australia","language":"en","unit":"m"},"location":{"name":"Sydney","country":"Australia","region":"New South Wales","lat":"-33.883","lon":"151.217","timezone_id":"Australia\/Sydney","localtime":"2022-03-14 17:20","localtime_epoch":1647278400,"utc_offset":"11.0"},"current":{"observation_time":"06:20 AM","temperature":22,"weather_code":116,"weather_icons":["https:\/\/assets.weatherstack.com\/images\/wsymbols01_png_64\/wsymbol_0002_sunny_intervals.png"],"weather_descriptions":["Partly cloudy"],"wind_speed":15,"wind_degree":170,"wind_dir":"S","pressure":1020,"precip":0,"humidity":78,"cloudcover":75,"feelslike":25,"uv_index":6,"visibility":10,"is_day":"yes"}}`

	t.Run("when no query string is provided", func(t *testing.T) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(standard_response)))
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
		}
	})
}
