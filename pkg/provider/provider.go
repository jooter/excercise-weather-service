// Implement Failsafe Weather Provider, which support failover and caching.
// The cache expire in 3 seconds.
// In other words, we will hit external provider at most once in 3 seconds.
package provider

import (
	"log"
	"sync"
	"time"

	"github.com/jooter/exercise-weather-service/pkg/core"
)

type weatherCache struct {
	weather    *core.Weather
	refreshied int64
	sync.RWMutex
}

// Failsafe Weather Provider
type FailsafeProvider struct {
	providerA   core.WeatherProvider
	providerB   core.WeatherProvider
	cache       *weatherCache
	cacheExpiry int // unit is second
}

func NewFailsafeProvider(providerA core.WeatherProvider, providerB core.WeatherProvider, expiry int) *FailsafeProvider {
	log.Println("A instance of failsafe provider has been create.")
	cache := &weatherCache{}
	return &FailsafeProvider{providerA: providerA, providerB: providerB, cache: cache, cacheExpiry: expiry}
}

// Implement interface core.WeatherProvider
func (p *FailsafeProvider) GetWeather() (*core.Weather, error) {
	start := time.Now().Unix()

	// it may need to wait, in case the cache is refreshing at this moment
	p.cache.RLock()

	lastRefreshed := p.cache.refreshied
	if p.cache.weather != nil {
		if start-lastRefreshed <= int64(p.cacheExpiry) {
			// cache is refreshed in no more than 3 seconds
			p.cache.RUnlock()
			return p.cache.weather, nil
		}
	}
	p.cache.RUnlock()

	// cache will be refreshed now

	p.cache.Lock()
	defer p.cache.Unlock()

	weather, err := p.providerA.GetWeather()
	if err != nil {
		weather, err = p.providerB.GetWeather()
		if err != nil {
			if p.cache.weather == nil {
				weather = &core.Weather{Error: core.ErrorGetWeatherFromExternal.Error()}
				return weather, core.ErrorGetWeatherFromExternal
			}

			// both provider failed, whatever in cache will be returned
			err = nil
			weather = p.cache.weather
			return weather, err
		}
	}

	// refreshed successfully
	p.cache.weather = weather
	p.cache.refreshied = time.Now().Unix()

	return weather, err
}
