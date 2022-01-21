// Implement Advanced Weather Provider, which support failover and caching.
// The cache expire in 3 seconds.
// To be confirm: do we timeout limit for a provider? I will assume no to start with.
//
// If a external provider takes 3 seconds to retrieve data, and in that 3 seconds there are many request received, we will only hit external provider only once.
// In other words, we will hit external provider at most once in 3 seconds.
package provider

import (
	"sync"
	"time"

	"github.com/jooter/exercise-weather-service/pkg/core"
)

type weatherCache struct {
	weather    *core.Weather
	refreshied int64
	sync.RWMutex
}

// Advanced Weather Provider
type FailsafeProvider struct {
	providerA core.WeatherProvider
	providerB core.WeatherProvider
	cache     *weatherCache
}

func NewFailsafeProvider(providerA core.WeatherProvider, providerB core.WeatherProvider) *FailsafeProvider {
	cache := &weatherCache{}
	return &FailsafeProvider{providerA: providerA, providerB: providerB, cache: cache}
}

func (p *FailsafeProvider) GetWeather() (*core.Weather, error) {
	start := time.Now().Unix()

	// it may need to wait, in case the cache is refreshing at this moment
	p.cache.RLock()

	lastRefreshed := p.cache.refreshied
	if p.cache.weather != nil {
		if start-lastRefreshed <= 3 {
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
				return nil, core.ErrorGetWeatherFromExternal
			}
			err = nil
			weather = p.cache.weather // both provider failed, whatever in cache will be returned
		}
	} else {
		// refreshed successfully
		p.cache.weather = weather
		p.cache.refreshied = time.Now().Unix()
	}
	return weather, err
}
