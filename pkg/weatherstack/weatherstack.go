// Implement WeatherProvider for using external weatherstack service
package weatherstack

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jooter/exercise-weather-service/pkg/core"
)

type Weatherstack struct {
	URL string
}

func NewWeatherstack(accessKey string) *Weatherstack {
	// skipped data validation
	return &Weatherstack{URL: "http://api.weatherstack.com/current?query=Melbourne&access_key=" + accessKey}
}

func (w Weatherstack) GetWeather() (*core.Weather, error) {

	ws, err := w.getWeatherstackResponse()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &core.Weather{TemperatureDegrees: ws.Current.Temperature, WindSpeed: ws.Current.WindSpeed}, nil
}

func (w Weatherstack) getWeatherstackResponse() (ws *weatherStackResponse, err error) {
	log.Println("connect:", w.URL)
	resp, err := http.Get(w.URL)
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
