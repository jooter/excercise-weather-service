# exercise-weather-service

## Requirement and discussion

* The service can hard-code Melbourne as a city.

    I assume supporting multiple cities/locations is out of scope of this exercise. 

    If query city is not Melbourne, "unknown city" error will be returned.

    If multi-cities is required, we can implement with redis db or in-memory cache as a map with RWMutex on map keys and RWMutex in each element of the map.

* The service should return a JSON payload with a unified response containing temperature in degrees Celsius and wind speed.

    I find wind speed from these two sites are very different. Please check [issues.md](issues.md) for details.

    When error, none OK status code will be returned, and playload like below will be returned.
    ```
    {
    	"error": "query parameter city is missing"
    }

    ```

* If one of the providers goes down, your service can quickly failover to a different provider without affecting your customers.

    For "quickly failover", configurable timeout has been implemented. [config.json](config.json) is an example config file.

* Have scalability and reliability in mind when designing the solution.
    - Docker or serverless could be used for scale out (not implement).
    - Query rate limit middle-ware could be implement to improve reliability (not implement).
    - Many errors and corner cases have investigated and tested.

* Weather results are fine to be cached for up to 3 seconds on the server in normal behaviour to prevent hitting weather providers.

    For implementing this, http handler is created with an instance of cache, and the workflow is:

    - Get read lock for cache
    - Check if cache is younger than 3 seconds (inclusive), if yes return cache. Otherwise, next step
    - Release read lock
    - Get read/write lock
    - Query on external provider and refresh cache
    - Release read/write lock

* Cached results should be served if all weather providers are down.

    Additionally, if no cache has been created yet, and all providers are down, an error message will be returned.
    ```
    {
	    "error": "error on getting weather from external"
    }
    ```

* The proposed solution should allow new developers to make changes to the code safely.

    For achieve this, I have been trying to:
    
    - Follow all SOLID principles
        - For example, dependency inversion principle has been used for designing failsafe provider in this solution.
        - The most of modules are assembled in main.go file by dependency injection technique.
    - Follow TDD process and the important modules have reached more than 90% coverage
    - Write "Clean Code" 

* Due to size and purpose of this project, compare with normal project I have skipped or been relexed for below:
    - My integration tests results are not verified automatically.
    - Docker container could be used for scaling out or dev/test for micro services.
    - Web service frameworks could be used for simplifying and standardize micro service implementation.
    - I also skipped unit tests for some module.
    - Access keys should not be keeping in source code repository.

## Build and test in Linux

In command line under the directory of this project
```
make
```
or 
```
make utest
make itest
```
or check build and test history in github actions tab in this project