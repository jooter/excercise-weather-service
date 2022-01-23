// Parse config file for web server
// this module could be re-implement with a popular lib spf13/viper for more complicated requirement.
package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	KeyOfWeatherstack   string `json:"key_of_weatherstack"`
	KeyOfOpenweathermap string `json:"key_of_openweathermap"`
	Timeout             int    `json:"timeout"`      // unit is second
	CacheExpiry         int    `json:"cache_expiry"` // unit is second
}

func ParseConfigFile(filename string) *Config {
	// set default values
	conf := Config{Timeout: 15, CacheExpiry: 3}
	configFile, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil
	}
	bs, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Println()
		return nil
	}
	err = json.Unmarshal(bs, &conf)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &conf
}
