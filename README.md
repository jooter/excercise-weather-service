# exercise-weather-service

## Requirement and discussion

* The service can hard-code Melbourne as a city.

    I assume supporting multiple cities/locations is out of scope of this exercise. 

    If multi-cities is required, we can implement with redis db or in-memory cache as a map with RWMutex with the whole map and RWMutex in each element of the map.

* The service should return a JSON payload with a unified response containing temperature in degrees Celsius and wind speed.

    I find wind speed from these two sites are very different. Please check [issues.md](issues.md) for details.

* If one of the providers goes down, your service can quickly failover to a different provider without affecting your customers.

    For "quickly failover", configurable timeout has been implemented.

* Have scalability and reliability in mind when designing the solution.

* Weather results are fine to be cached for up to 3 seconds on the server in normal behaviour to prevent hitting weather providers.

    For implementing this, http handler is created with a instance of cache, and the workflow is:

    - Get read lock for cache
    - Check if cache is younger than 3 seconds (inclusive), if yes return cache. Otherwise next step
    - Release read lock
    - Get read/write lock
    - Query on external provider and refresh cache
    - Release read/write lock

* Cached results should be served if all weather providers are down.

    If no cache has been created, and all providers are down, error message will be returned.

* The proposed solution should allow new developers to make changes to the code safely.

    For achieve this, I have been trying:
    
    - Follow SOLID principles