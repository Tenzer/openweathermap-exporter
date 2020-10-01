# OpenWeatherMap Exporter

[Prometheus](https://prometheus.io/) exporter for the [OpenWeatherMap](https://openweathermap.org/) service.

This requires you to sign up on their website in order to receive a free API key.


## Metrics

| Name | Description |
|---|---|
| `weather_cloudiness_percent` | Cloud cover in percentage |
| `weather_feelslike_celsius` | Current temperature taking human perception into account, in celsius |
| `weather_humidity_percent` | Humidity in percent |
| `weather_pressure_hpa` | Atmospheric pressure in hectopascal |
| `weather_sunrise_timestamp` | Sunrise time as a UNIX timestamp |
| `weather_sunset_timestamp` | Sunset time as a UNIX timestamp |
| `weather_temperature_celsius` | Current temperature in celsius |
| `weather_temperature_max_celsius` | Current maximum temperature in celsius |
| `weather_temperature_min_celsius` | Current minimum temperature in celsius |
| `weather_up` | Whether the metrics can be collected |
| `weather_visibility_meters` | Visibility in meters |
| `weather_winddirection_degrees` | Wind direction in degrees |
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
