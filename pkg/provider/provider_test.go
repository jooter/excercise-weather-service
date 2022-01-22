package provider

import (
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/jooter/exercise-weather-service/pkg/core"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

type mockProvider struct {
	windSpeed float32
	err       error
	sleep     int
	counter   int // how many times of refresh
}

func newMockProvider(windSpeed float32, err error, sleep int) *mockProvider {
	return &mockProvider{windSpeed: windSpeed, err: err, sleep: sleep}
}

func (p *mockProvider) GetWeather() (*core.Weather, error) {
	log.Println("the cache will be refreshed now")
	time.Sleep(time.Duration(p.sleep) * time.Second)
	p.counter++
	return &core.Weather{WindSpeed: p.windSpeed}, p.err
}

func TestBasic(t *testing.T) {
	mock01 := newMockProvider(1, nil, 0)
	mock02 := newMockProvider(2, nil, 0)
	pvd := NewFailsafeProvider(mock01, mock02, 3)
	w, err := pvd.GetWeather()
	assert.NoError(t, err)
	assert.Equal(t, float32(1), w.WindSpeed)
}

func TestFailover(t *testing.T) {
	mock01 := newMockProvider(1.1, errors.New("test error"), 0)
	mock02 := newMockProvider(2, nil, 0)
	pvd := NewFailsafeProvider(mock01, mock02, 3)
	w, err := pvd.GetWeather()
	assert.NoError(t, err)
	assert.Equal(t, float32(2), w.WindSpeed)
}

func TestOnlyCacheToUse(t *testing.T) {
	mock01 := newMockProvider(1.1, errors.New("test error"), 0)
	mock02 := newMockProvider(2, errors.New("test error 2"), 0)
	pvd := NewFailsafeProvider(mock01, mock02, 3)
	pvd.cache.weather = &core.Weather{}
	w, err := pvd.GetWeather()
	assert.NoError(t, err)
	assert.Equal(t, float32(0), w.WindSpeed)
}

func TestError(t *testing.T) {
	mock01 := newMockProvider(1.1, errors.New("test error"), 0)
	mock02 := newMockProvider(2, errors.New("test error 2"), 0)
	pvd := NewFailsafeProvider(mock01, mock02, 3)
	_, err := pvd.GetWeather()
	assert.Error(t, err)
}

func TestSequentialRequest(t *testing.T) {
	mock01 := newMockProvider(1, nil, 2)
	pvd := NewFailsafeProvider(mock01, mock01, 3)
	start := time.Now().Unix()
	log.Println("start test")
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		w, err := pvd.GetWeather()
		log.Println("weather,err =", w, err)
	}
	log.Println("test took", time.Now().Unix()-start, "seconds")
	assert.True(t, mock01.counter == 2)
}

func TestParallelRequest(t *testing.T) {
	mock01 := newMockProvider(1, nil, 1)
	pvd := NewFailsafeProvider(mock01, mock01, 3)
	start := time.Now().Unix()
	log.Println("start test")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			w, err := pvd.GetWeather()
			log.Println("weather,err =", w, err)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println("test took", time.Now().Unix()-start, "seconds")
	assert.True(t, mock01.counter == 1)
}
