# OpenWeatherMap Exporter

[Prometheus](https://prometheus.io/) exporter for the [OpenWeatherMap](https://openweathermap.org/) service.

This requires you to sign up on their website in order to receive a free API key.


## Metrics

| Name | Description |
|---|---|
| `weather_cloudiness_percent` | Cloud cover in percentage |
| `weather_feelslike_celsius` | Current temperature taking human perception into account, in celsius |
| `weather_feelslike_fahrenheit` | Current temperature taking human perception into account, in fahrenheit |
| `weather_feelslike_kelvin` | Current temperature taking human perception into account, in kelvin |
| `weather_humidity_percent` | Humidity in percent |
| `weather_maxtemperature_celsius` | Current maximum temperature in celsius |
| `weather_maxtemperature_fahrenheit` | Current maximum temperature in fahrenheit |
| `weather_maxtemperature_kelvin` | Current maximum temperature in kelvin |
| `weather_mintemperature_celsius` | Current minimum temperature in celsius |
| `weather_mintemperature_fahrenheit` | Current minimum temperature in fahrenheit |
| `weather_mintemperature_kelvin` | Current minimum temperature in kelvin |
| `weather_pressure_hpa` | Atmospheric pressure in hectopascal |
| `weather_sunrise` | Sunrise time as a UNIX timestamp |
| `weather_sunset` | Sunset time as a UNIX timestamp |
| `weather_temperature_celsius` | Current temperature in celsius |
| `weather_temperature_fahrenheit` | Current temperature in fahrenheit |
| `weather_temperature_kelvin` | Current temperature in kelvin |
| `weather_up` | Whether the metrics can be collected |
| `weather_visibility_meters` | Visibility in meters |
| `weather_visibility_miles` | Visibility in miles |
| `weather_winddirection` | Wind direction in degrees |
| `weather_windspeed_mph` | Wind speed in miles per hour |
| `weather_windspeed_mps` | Wind speed in meters per second |

All metrics with the exception of `weather_up` has the following labels with information about the location the metric is reported for:

| Name | Description |
|---|---|
| `location_id` | OpenWeatherMap location ID |
| `location_name` | Name of the location |
| `location_country` | 2 letter country code |
| `latitude` | Latitude |
| `longitude` | Longitude |


## Building and running

If you have [Go](https://golang.org/) installed you can download and build the exporter with:

    go get github.com/Tenzer/openweathermap-exporter

The following configuration flags are available:

    $ openweathermap-exporter -help
    Usage of openweathermap-exporter:
      -apikey string
            API key for OpenWeatherMap (required)
      -listen-address string
            Hostname and port to listen on (default "localhost:9755")
      -location-ids string
            Comma separated list of location IDs to fetch weather for (required)

You can look up the location IDs on the [OpenWeatherMap website](https://openweathermap.org/find).


## Note

The exporter makes an API request to OpenWeatherMap every time the `/metrics` endpoint is requested.
This allows you to control how often the data is fetched from them.
Be aware that they [specifically say](https://openweathermap.org/appid#work) that weather data isn't updated more often than every 10 minutes,
so it likely won't make sense to request it more often than that.
