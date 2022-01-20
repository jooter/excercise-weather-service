package openweathermap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWeatherError01(t *testing.T) {
	ws := NewWeatherstack("")
	w, err := ws.GetWeather()
	assert.NotNil(t, err)
	assert.Empty(t, w)
}
