package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWeather(t *testing.T) {
	t.Run("when no query string is provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
		res := httptest.NewRecorder()

		GetWeather(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.Code)
		}
	})
}
