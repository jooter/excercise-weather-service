package openweathermap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ws := New("")
	w, err := ws.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, w)
}

func TestGetWeatherOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// dummy response wind speed 1 m/s ( 3.6 km/h)
		fmt.Fprintln(w, `{
			"wind": {
				"speed": 1
			}
		}`)
	}))
	client := &Openweathermap{URL: srv.URL}
	weather, err := client.GetWeather()
	assert.Nil(t, err)
	assert.Equal(t, float32(3.6), weather.WindSpeed)
}

func TestGetWeatherNetworkError(t *testing.T) {
	client := &Openweathermap{URL: "wrong://"}
	weather, err := client.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, weather)
}
func TestGetWeatherTokenError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// dummy response wind speed 1 m/s ( 3.6 km/h)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `{"message": "tokern error"}`)
	}))
	client := &Openweathermap{URL: srv.URL}
	weather, err := client.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, weather)
}

func TestGetWeatherErrorOnJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// dummy response wind speed 1 m/s ( 3.6 km/h)
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, `not json here`)
	}))
	client := &Openweathermap{URL: srv.URL}
	weather, err := client.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, weather)
}
