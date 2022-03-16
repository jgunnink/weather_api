package api

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jgunnink/weather_api/utils/mocks"
)

func init() {
	Client = &mocks.MockClient{}
}

func TestGetWeather(t *testing.T) {
	Client = &mocks.MockClient{}
	weatherstack_api_response := `{"request":{"type":"City","query":"Sydney, Australia","language":"en","unit":"m"},"location":{"name":"Sydney","country":"Australia","region":"New South Wales","lat":"-33.883","lon":"151.217","timezone_id":"Australia\/Sydney","localtime":"2022-03-14 17:20","localtime_epoch":1647278400,"utc_offset":"11.0"},"current":{"observation_time":"06:20 AM","temperature":22,"weather_code":116,"weather_icons":["https:\/\/assets.weatherstack.com\/images\/wsymbols01_png_64\/wsymbol_0002_sunny_intervals.png"],"weather_descriptions":["Partly cloudy"],"wind_speed":15,"wind_degree":170,"wind_dir":"S","pressure":1020,"precip":0,"humidity":78,"cloudcover":75,"feelslike":25,"uv_index":6,"visibility":10,"is_day":"yes"}}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(weatherstack_api_response)))

	t.Run("when no query string is provided, it defaults to Sydney", func(t *testing.T) {
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

	t.Run("the response includes just the wind speed and temperature", func(t *testing.T) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(weatherstack_api_response)))
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		expectedResponse := `{"wind_speed":15,"temperature_degrees":22}`
		actualResponse := strings.Replace(res.Body.String(), "\n", "", -1) // remove newlines to make comparison easy
		if actualResponse != expectedResponse {
			t.Errorf("Expected body %q, got %q", expectedResponse, actualResponse)
		}
	})

	t.Run("When there is a query string provided, it uses the query", func(t *testing.T) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(weatherstack_api_response)))
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather?query=Perth", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
		}
	})

	t.Run("When an unknown query string is provided, it defaults to Sydney", func(t *testing.T) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(weatherstack_api_response)))
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather?potato=Sydney", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
		}
	})

	t.Run("When there is an error, it returns a 503", func(t *testing.T) {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return nil, errors.New("error")
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, res.Code)
		}
	})

	t.Run("open weathermap is used when weatherstack is offline", func(t *testing.T) {
		t.Skip() // TODO: Currently broken. Mocking the getdofunc to return a 500 also fails the second request.
		// Need to work out how to mock the request once, and the subsequent request to return a 200 with a sample payload.
		// from the OpenWeatherMap API.
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather?query=Sydney", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
		}
	})
}
