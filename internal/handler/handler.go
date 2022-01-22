package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jooter/exercise-weather-service/pkg/core"
	"github.com/jooter/exercise-weather-service/pkg/provider"
)

func CreateWeatherHandler(failsafeProvider *provider.FailsafeProvider) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		weather, err := failsafeProvider.GetWeather()
		if err != nil && err != core.ErrorGetWeatherFromExternal {
			log.Println(err)
			return
		}
		b, err := json.MarshalIndent(weather, "", "\t")
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Fprintln(w, string(b))
	}
}
