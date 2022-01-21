// Implement WeatherProvider for using external weatherstack service
package weatherstack

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jooter/exercise-weather-service/pkg/core"
)

type Weatherstack struct {
	url     string
	timeout int
}

func New(accessKey string, timeout int) *Weatherstack {
	// skipped data validation
	return &Weatherstack{timeout: timeout,
		url: "http://api.weatherstack.com/current?query=Melbourne&access_key=" + accessKey}
}

func (w Weatherstack) GetWeather() (*core.Weather, error) {

	ws, err := w.request()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &core.Weather{TemperatureDegrees: ws.Current.Temperature, WindSpeed: ws.Current.WindSpeed}, nil
}

func (w Weatherstack) request() (ws *weatherStackResponse, err error) {
	log.Println("connect:", w.url)
	client := http.Client{Timeout: time.Duration(w.timeout) * time.Second}
	resp, err := client.Get(w.url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println("status code =", resp.StatusCode)
		return nil, core.ErrorGetWeatherFromExternal
	}

	err = dec.Decode(&ws)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	errInfo := ws.Error.Info
	if errInfo != "" {
		log.Println("external error :", errInfo)
		return nil, core.ErrorGetWeatherFromExternal
	}

	return
}
