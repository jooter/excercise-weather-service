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
	Timeout             int    `json:"timeout"`
	CacheExpiry         int    `json:"cache_expiry"`
}

func ParseConfigFile(filename string) *Config {
	var conf Config
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
