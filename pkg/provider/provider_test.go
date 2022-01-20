package provider

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/jooter/exercise-weather-service/pkg/core"
	"github.com/stretchr/testify/assert"
)

type mockProvider struct {
	windSpeed float32
	err       error
	sleep     int
}

func newMockProvider(windSpeed float32, err error, sleep int) *mockProvider {
	return &mockProvider{windSpeed: windSpeed, err: err, sleep: sleep}
}

func (p mockProvider) GetWeather() (*core.Weather, error) {
	time.Sleep(time.Duration(p.sleep) * time.Second)
	return &core.Weather{WindSpeed: p.windSpeed}, p.err
}

func TestBasic(t *testing.T) {
	mock01 := newMockProvider(1, nil, 0)
	mock02 := newMockProvider(2, nil, 0)
	adv := NewAdvancedProvider(mock01, mock02)
	w, err := adv.GetWeather()
	assert.NoError(t, err)
	assert.Equal(t, float32(1), w.WindSpeed)
}

func TestFailover(t *testing.T) {
	mock01 := newMockProvider(1.1, errors.New("test error"), 0)
	mock02 := newMockProvider(2, nil, 0)
	adv := NewAdvancedProvider(mock01, mock02)
	w, err := adv.GetWeather()
	assert.NoError(t, err)
	assert.Equal(t, float32(2), w.WindSpeed)
}

func TestOnlyCacheToUse(t *testing.T) {
	mock01 := newMockProvider(1.1, errors.New("test error"), 0)
	mock02 := newMockProvider(2, errors.New("test error 2"), 0)
	adv := NewAdvancedProvider(mock01, mock02)
	w, err := adv.GetWeather()
	assert.NoError(t, err)
	assert.Equal(t, float32(0), w.WindSpeed)
}

func TestSequentialRequest(t *testing.T) {
	mock01 := newMockProvider(1, nil, 10)
	adv := NewAdvancedProvider(mock01, mock01)
	start := time.Now().Unix()
	log.Println("start")
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		w, err := adv.GetWeather()
		log.Println(w, err)
	}
	log.Println(time.Now().Unix() - start)
	assert.Equal(t, int64(40), time.Now().Unix()-start,
		"Total time should be 40 sec (cache is refreshed 3 times in 10 sec, each time takes 10 sec; 10 requests wait 10 sec)")
}

func TestParallelRequest(t *testing.T) {
	mock01 := newMockProvider(1, nil, 10)
	adv := NewAdvancedProvider(mock01, mock01)
	start := time.Now().Unix()
	log.Println("start")
	for i := 0; i < 10; i++ {
		go func() {
			w, err := adv.GetWeather()
			log.Println(w, err)
		}()
	}
	adv.GetWeather()
	log.Println(time.Now().Unix() - start)
	assert.Equal(t, int64(10), time.Now().Unix()-start,
		"Total time should be 10 sec. Only one of request will hit external provider, which takes 10 sec.")
}
