package handler

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jooter/exercise-weather-service/pkg/core"
	"github.com/jooter/exercise-weather-service/pkg/provider"
	"github.com/stretchr/testify/assert"
)

func TestNotFoundHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	NotFoundHandler(w, r)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	expect := `{
	"error": "not found"
}
`
	assert.Equal(t, expect, string(data))
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestQueryErrorMissingParameter(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	hdl := CreateWeatherHandler(nil)
	hdl.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	expect := `{
	"error": "query parameter city is missing"
}
`
	assert.Equal(t, expect, string(data))
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestQueryErrorUnknownCity(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/?city=sydney", nil)
	hdl := CreateWeatherHandler(nil)
	hdl.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	expect := `{
	"error": "unknown city"
}
`
	assert.Equal(t, expect, string(data))
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
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

func TestFailsafeHandlerOK(t *testing.T) {
	mock01 := newMockProvider(1, nil, 0)
	pvd := provider.NewFailsafeProvider(mock01, mock01, 3)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/?city=Melbourne", nil)
	hdl := CreateWeatherHandler(pvd)
	hdl.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	expect := `{
	"wind_speed": 1
}
`
	assert.Equal(t, expect, string(data))
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestFailsafeHandlerError(t *testing.T) {
	mock01 := newMockProvider(1, errors.New("test error"), 0)
	pvd := provider.NewFailsafeProvider(mock01, mock01, 3)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/?city=Melbourne", nil)
	hdl := CreateWeatherHandler(pvd)
	hdl.ServeHTTP(w, r)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	expect := `{
	"error": "error on getting weather from external"
}
`
	assert.Equal(t, expect, string(data))
	assert.Equal(t, http.StatusFailedDependency, res.StatusCode)
}
