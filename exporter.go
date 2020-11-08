package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels = []string{"location_id", "location_name", "location_country", "latitude", "longitude"}

	cloudiness     = prometheus.NewDesc("weather_cloudiness_percent", "Cloud cover in percentage", labels, nil)
	feelslike      = prometheus.NewDesc("weather_feelslike_celsius", "Current temperature taking human perception into account, in celsius", labels, nil)
	humidity       = prometheus.NewDesc("weather_humidity_percent", "Humidity in percent", labels, nil)
	pressure       = prometheus.NewDesc("weather_pressure_hpa", "Atmospheric pressure in hectopascal", labels, nil)
	sunrise        = prometheus.NewDesc("weather_sunrise_timestamp", "Sunrise time as a UNIX timestamp", labels, nil)
	sunset         = prometheus.NewDesc("weather_sunset_timestamp", "Sunset time as a UNIX timestamp", labels, nil)
	temperature    = prometheus.NewDesc("weather_temperature_celsius", "Current temperature in celsius", labels, nil)
	temperatureMax = prometheus.NewDesc("weather_temperature_max_celsius", "Current maximum temperature in celsius", labels, nil)
	temperatureMin = prometheus.NewDesc("weather_temperature_min_celsius", "Current minimum temperature in celsius", labels, nil)
	up             = prometheus.NewDesc("weather_up", "Whether the metrics can be collected", nil, nil)
	visibility     = prometheus.NewDesc("weather_visibility_meters", "Visibility in meters", labels, nil)
	windDirection  = prometheus.NewDesc("weather_winddirection_degrees", "Wind direction in degrees", labels, nil)
	windSpeed      = prometheus.NewDesc("weather_windspeed_mps", "Wind speed in meters per second", labels, nil)
)

type Exporter struct {
	apikey          string
	cityIds         string
	cacheTTL        time.Duration
	cachedResponse  response
	responseFetched time.Time
}

func newExporter(apikey string, cityIds string, cacheTTL time.Duration) *Exporter {
	return &Exporter{
		apikey:   apikey,
		cityIds:  cityIds,
		cacheTTL: cacheTTL,
	}
}

func (exporter *Exporter) Describe(channel chan<- *prometheus.Desc) {
	channel <- cloudiness
	channel <- feelslike
	channel <- humidity
	channel <- pressure
	channel <- sunrise
	channel <- sunset
	channel <- temperature
	channel <- temperatureMax
	channel <- temperatureMin
	channel <- up
	channel <- visibility
	channel <- windDirection
	channel <- windSpeed
}

func (exporter *Exporter) Collect(channel chan<- prometheus.Metric) {
	data := exporter.getData()
	if len(data.Locations) == 0 {
		channel <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 0)
		return
	}

	channel <- prometheus.MustNewConstMetric(up, prometheus.GaugeValue, 1)

	for _, location := range data.Locations {
		locationID := strconv.FormatInt(location.ID, 10)
		latitude := fmt.Sprintf("%f", location.Coordinates.Latitude)
		longitude := fmt.Sprintf("%f", location.Coordinates.Longitude)

		channel <- prometheus.MustNewConstMetric(cloudiness, prometheus.GaugeValue, location.Clouds.Cover, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(feelslike, prometheus.GaugeValue, location.Main.FeelsLike, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(humidity, prometheus.GaugeValue, location.Main.Humidity, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(pressure, prometheus.GaugeValue, location.Main.Pressure, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(sunrise, prometheus.GaugeValue, location.Sys.Sunrise, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(sunset, prometheus.GaugeValue, location.Sys.Sunset, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(temperature, prometheus.GaugeValue, location.Main.Temperature, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(temperatureMax, prometheus.GaugeValue, location.Main.TemperatureMax, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(temperatureMin, prometheus.GaugeValue, location.Main.TemperatureMin, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(visibility, prometheus.GaugeValue, location.Visibility, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(windDirection, prometheus.GaugeValue, location.Wind.Direction, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(windSpeed, prometheus.GaugeValue, location.Wind.Speed, locationID, location.Name, location.Sys.Country, latitude, longitude)
	}
}

func (exporter *Exporter) getData() response {
	if time.Since(exporter.responseFetched) < exporter.cacheTTL {
		return exporter.cachedResponse
	}

	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/group?units=metric&id=%s&appid=%s", exporter.cityIds, exporter.apikey))
	if err != nil {
		log.Println("Request to OpenWeatherMap failed:", err)
		return response{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Reading response body failed:", err)
		return response{}
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Response had unexpected status code:", resp.StatusCode, "- Response body:", string(body))
		return response{}
	}

	data := response{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Parsing JSON response failed:", err)
		return response{}
	}

	exporter.cachedResponse = data
	exporter.responseFetched = time.Now()

	return data
}
