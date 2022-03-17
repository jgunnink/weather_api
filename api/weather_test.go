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
	// openweathermap_api_response := `{"coord": {"lon": -122.08,"lat": 37.39},"weather": [{"id": 800,"main": "Clear","description": "clear sky","icon": "01d"}],"base": "stations","main": {"temp": 282.55,"feels_like": 281.86,"temp_min": 280.37,"temp_max": 284.26,"pressure": 1023,"humidity": 100},"visibility": 10000,"wind": {"speed": 1.5,"deg": 350},"clouds": {"all": 1},"dt": 1560350645,"sys": {"type": 1,"id": 5122,"message": 0.0139,"country": "US","sunrise": 1560343627,"sunset": 1560396563},"timezone": -25200,"id": 420006353,"name": "Mountain View","cod": 200}`

	r := ioutil.NopCloser(bytes.NewReader([]byte(weatherstack_api_response)))

	t.Run("When no query string is provided, it defaults to Sydney", func(t *testing.T) {
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

	t.Run("When the city cannot be found", func(t *testing.T) {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       r,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, res.Code)
		}
	})

	t.Run("When there are too many requests", func(t *testing.T) {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       r,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, res.Code)
		}
	})

	t.Run("When there is a timeout from the upstream", func(t *testing.T) {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusGatewayTimeout,
				Body:       nil,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, res.Code)
		}
	})

	t.Run("When there is some other status code", func(t *testing.T) {
		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusTeapot,
				Body:       nil,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, res.Code)
		}
	})

	t.Run("When WeatherStack responds with a custom error on a status 200", func(t *testing.T) {
		customResponse := `{"success":false,"error":{"code":615,"type":"request_failed","info":"Your API request failed. Please try again or contact support."}}`
		cr := ioutil.NopCloser(bytes.NewReader([]byte(customResponse)))

		mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       cr,
			}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/v1/weather?potato=Sydney", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, res.Code)
		}
	})

	// TODO: Currently broken. Mocking the getdofunc to return a 500 also fails the second request.
	// Need to work out how to mock the request once, and the subsequent request to return a 200 with a sample payload.
	// from the OpenWeatherMap API.
	t.Run("open weathermap is used when weatherstack is offline", func(t *testing.T) {
		t.Skip()
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
