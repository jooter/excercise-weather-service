package openweathermap

// Implement WeatherProvider for using external weatherstack service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jooter/exercise-weather-service/pkg/core"
)

type Openweathermap struct {
	URL string
}

func NewWeatherstack(accessKey string) *Openweathermap {
	// skipped data validation

	// Add units=metric in order to getting temperature in celsius
	return &Openweathermap{URL: "http://api.openweathermap.org/data/2.5/weather?q=melbourne,AU&units=metric&appid=" + accessKey}
}

func (w Openweathermap) GetWeather() (*core.Weather, error) {

	ws, err := w.getWeatherstackResponse()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// Wind speed unit is m/s from source. It convert to km/h by multiply 3.6 .
	return &core.Weather{TemperatureDegrees: ws.Main.Temp, WindSpeed: ws.Wind.Speed * 3.6}, nil
}

func (w Openweathermap) getWeatherstackResponse() (ws *openweathermapResponse, err error) {
	log.Println("connect:", w.URL)
	resp, err := http.Get(w.URL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ws)

	if resp.StatusCode != http.StatusOK || err != nil || ws.Message != "" {
		log.Println("status code =", resp.StatusCode)
		if err != nil {
			log.Println(err)
		}
		if ws.Message != "" {
			log.Println("external error :", ws.Message)
		}
		return nil, core.ErrorGetWeatherFromExternal
	}

	return
}
