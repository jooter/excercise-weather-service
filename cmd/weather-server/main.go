package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jooter/exercise-weather-service/internal/config"
	"github.com/jooter/exercise-weather-service/internal/handler"
	"github.com/jooter/exercise-weather-service/pkg/openweathermap"
	"github.com/jooter/exercise-weather-service/pkg/provider"
	"github.com/jooter/exercise-weather-service/pkg/weatherstack"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// Get configuration
	configFilenamePtr := flag.String("conf", "config.json", "JSON config file")
	serverAddressPtr := flag.String("addr", ":8080", "Server listening at address")
	flag.Parse()
	conf := config.ParseConfigFile(*configFilenamePtr)
	if conf == nil {
		log.Fatalln("config file parsing is failed")
	}

	// Prepare provider
	providerA := weatherstack.New(conf.KeyOfWeatherstack, conf.Timeout)
	providerB := openweathermap.New(conf.KeyOfOpenweathermap, conf.Timeout)
	failsafeProvider := provider.NewFailsafeProvider(providerA, providerB, conf.CacheExpiry)

	// HTTP server
	weatherHandler := handler.CreateWeatherHandler(failsafeProvider)
	http.Handle("/v1/weather", weatherHandler)
	http.HandleFunc("/", handler.NotFoundHandler)

	log.Println("Listen at", *serverAddressPtr)
	log.Fatalln(http.ListenAndServe(*serverAddressPtr, nil))
	// graceful shutdown has been skipped
}
