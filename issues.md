# Issues to be investigate

## Wind speed from two web api are very different
```
// from weatherstack
{
    "request": {
        "type": "City",
        "query": "Melbourne, Australia",
        "language": "en",
        "unit": "m"
    },
    "location": {
        "name": "Melbourne",
        "country": "Australia",
        "region": "Victoria",
        "lat": "-37.817",
        "lon": "144.967",
        "timezone_id": "Australia/Melbourne",
        "localtime": "2022-01-21 07:37",
        "localtime_epoch": 1642750620,
        "utc_offset": "11.0"
    },
    "current": {
        "observation_time": "08:37 PM",
        "temperature": 18,
        "weather_code": 113,
        "weather_icons": [
            "https://assets.weatherstack.com/images/wsymbols01_png_64/wsymbol_0001_sunny.png"
        ],
        "weather_descriptions": [
            "Sunny"
        ],
        "wind_speed": 20, // unit km/h
        "wind_degree": 10,
        "wind_dir": "N",
        "pressure": 1027,
        "precip": 0,
        "humidity": 64,
        "cloudcover": 0,
        "feelslike": 18,
        "uv_index": 1,
        "visibility": 10,
        "is_day": "yes"
    }
}

// from openweather map
{
    "coord": {
        "lon": 144.9633,
        "lat": -37.814
    },
    "weather": [
        {
            "id": 800,
            "main": "Clear",
            "description": "clear sky",
            "icon": "01d"
        }
    ],
    "base": "stations",
    "main": {
        "temp": 18.12,
        "feels_like": 17.87,
        "temp_min": 13.43,
        "temp_max": 20.18,
        "pressure": 1025,
        "humidity": 72
    },
    "visibility": 10000,
    "wind": {
        "speed": 0.89, // unit m/s, equal to 3.2 km/h
        "deg": 85,
        "gust": 1.79
    },
    "clouds": {
        "all": 0
    },
    "dt": 1642710662, // GMT: Thursday, 20 January 2022 20:31:02
    "sys": {
        "type": 2,
        "id": 2008797,
        "country": "AU",
        "sunrise": 1642706448,
        "sunset": 1642758085
    },
    "timezone": 39600,
    "id": 2158177,
    "name": "Melbourne",
    "cod": 200
}

```
