package weatherstack

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ws := New("", 15)
	w, err := ws.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, w)
}

func TestGetWeatherOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// dummy response wind speed 20 km/h
		fmt.Fprintln(w, `{
			"current": {
				"wind_speed": 20
			}
		}`)
	}))
	client := &Weatherstack{url: srv.URL}
	weather, err := client.GetWeather()
	assert.Nil(t, err)
	assert.Equal(t, float32(20), weather.WindSpeed)
}

func TestGetWeatherNetworkError(t *testing.T) {
	client := &Weatherstack{url: "wrong://"}
	weather, err := client.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, weather)
}
func TestGetWeatherTokenError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"error": {"info": "tokern error"}}`)
	}))
	client := &Weatherstack{url: srv.URL}
	weather, err := client.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, weather)
}

func TestGetWeatherErrorOnJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `not json here`)
	}))
	client := &Weatherstack{url: srv.URL}
	weather, err := client.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, weather)
}

func TestGetWeatherErrorOnStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	client := &Weatherstack{url: srv.URL}
	weather, err := client.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, weather)
}
