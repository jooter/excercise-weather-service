// this package contents core struct and interface
package core

import "errors"

type Weather struct {
	TemperatureDegrees float32 `json:"temperature_degrees"`
	WindSpeed          float32 `json:"wind_speed"`
}

type WeatherProvider interface {
	GetWeather() (*Weather, error)
}

var (
	ErrorGetWeatherFromExternal error = errors.New("error on getting weather from external")
)
