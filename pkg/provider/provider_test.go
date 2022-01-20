package provider

import (
	"errors"
	"testing"

	"github.com/jooter/exercise-weather-service/pkg/core"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct {
	windSpeed float32
	err       error
}

func newMockProvider(windSpeed float32, err error) *mockProvider {
	return &mockProvider{windSpeed: windSpeed, err: err}
}

func (p mockProvider) GetWeather() (*core.Weather, error) {
	return &core.Weather{WindSpeed: p.windSpeed}, p.err
}

func TestBasic(t *testing.T) {
	mock01 := newMockProvider(1, nil)
	mock02 := newMockProvider(2, nil)
	adv := NewAdvancedProvider(mock01, mock02)
	w, err := adv.GetWeather()
	assert.NoError(t, err)
	assert.Equal(t, float32(1), w.WindSpeed)
}

func TestFailover(t *testing.T) {
	mock01 := newMockProvider(1.1, errors.New("test error"))
	mock02 := newMockProvider(2, nil)
	adv := NewAdvancedProvider(mock01, mock02)
	w, err := adv.GetWeather()
	assert.NoError(t, err)
	assert.Equal(t, float32(2), w.WindSpeed)
}
