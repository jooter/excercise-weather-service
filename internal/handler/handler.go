package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jooter/exercise-weather-service/pkg/core"
	"github.com/jooter/exercise-weather-service/pkg/provider"
)

func CreateWeatherHandler(failsafeProvider *provider.FailsafeProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !r.URL.Query().Has("city") {
			response(w, nil, http.StatusBadRequest, "query parameter city is missing")
			return
		}
		if strings.ToLower(r.URL.Query().Get("city")) != "melbourne" {
			response(w, nil, http.StatusBadRequest, "unknown city")
			return
		}
		weather, err := failsafeProvider.GetWeather()
		if err != nil {
			log.Println(err)
			response(w, nil, http.StatusFailedDependency, err.Error())
			return
		}
		response(w, weather, http.StatusOK, "")
	})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response(w, nil, http.StatusNotFound, "not found")
}

func response(w http.ResponseWriter, weather *core.Weather, code int, errStr string) {
	if weather == nil {
		weather = &core.Weather{Error: errStr}
	}

	w.WriteHeader(code)

	b, err := json.MarshalIndent(weather, "", "\t")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintln(w, string(b))
}
