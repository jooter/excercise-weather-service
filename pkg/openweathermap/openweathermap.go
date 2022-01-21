// Implement WeatherProvider for using external weatherstack service
package openweathermap

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jooter/exercise-weather-service/pkg/core"
)

type Openweathermap struct {
	url     string
	timeout int
}

func New(accessKey string, timeout int) *Openweathermap {
	// skipped data validation

	// Add units=metric in order to getting temperature in celsius
	return &Openweathermap{timeout: timeout,
		url: "http://api.openweathermap.org/data/2.5/weather?q=melbourne,AU&units=metric&appid=" + accessKey}
}

func (w Openweathermap) GetWeather() (*core.Weather, error) {

	ws, err := w.request()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Wind speed unit is m/s from source. It convert to km/h by multiply 3.6 .
	return &core.Weather{TemperatureDegrees: ws.Main.Temp, WindSpeed: ws.Wind.Speed * 3.6}, nil
}

func (w Openweathermap) request() (ws *openweathermapResponse, err error) {
	log.Println("connect:", w.url) // to be removed for provent key leaking
	client := http.Client{Timeout: time.Duration(w.timeout) * time.Second}
	resp, err := client.Get(w.url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ws)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK || ws.Message != "" {
		log.Println("status code =", resp.StatusCode)
		if ws.Message != "" {
			log.Println("external error :", ws.Message)
		}
		return nil, core.ErrorGetWeatherFromExternal
	}

	return
}
