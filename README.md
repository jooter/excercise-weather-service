# exercise-weather-service

## Requirement and discussion

* The service can hard-code Melbourne as a city.

I assume supporting multiple cities/locations is out of scope of this exercise. 

* The service should return a JSON payload with a unified response containing temperature in degrees Celsius and wind speed.

I find wind speed from these two sites are very different. Please check [issues.md](issues.md) for details.

* If one of the providers goes down, your service can quickly failover to a different provider without affecting your customers.

For "quickly failover", configurable timeout has been implemented.

* Have scalability and reliability in mind when designing the solution.

* Weather results are fine to be cached for up to 3 seconds on the server in normal behaviour to prevent hitting weather providers.

For implementing this, I create a instance of cache, and share access to all handler calls, and then the workflow is:

    - Check and wait cache refreshing
    - Check if cache is younger than 3 seconds (inclusive), return cache. Otherwise next step
    - Lock cache by adding 1 to wait group in cache struct
    - Query on external provider and refresh cache
    - Release lock by marking wait group as done

* Cached results should be served if all weather providers are down.

If no cache has been created, and all providers are down, error message will be returned.

* The proposed solution should allow new developers to make changes to the code safely.

For achieve this, I have done:
    
    - Follow SOLID principles