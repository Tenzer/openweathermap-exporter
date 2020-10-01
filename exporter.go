package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels = []string{"location_id", "location_name", "location_country", "latitude", "longitude"}

	cloudinessPercent        = prometheus.NewDesc("weather_cloudiness_percent", "Cloud cover in percentage", labels, nil)
	feelslikeCelsius         = prometheus.NewDesc("weather_feelslike_celsius", "Current temperature taking human perception into account, in celsius", labels, nil)
	feelslikeFahrenheit      = prometheus.NewDesc("weather_feelslike_fahrenheit", "Current temperature taking human perception into account, in fahrenheit", labels, nil)
	feelslikeKelvin          = prometheus.NewDesc("weather_feelslike_kelvin", "Current temperature taking human perception into account, in kelvin", labels, nil)
	humidityPercent          = prometheus.NewDesc("weather_humidity_percent", "Humidity in percent", labels, nil)
	maxTemperatureCelsius    = prometheus.NewDesc("weather_maxtemperature_celsius", "Current maximum temperature in celsius", labels, nil)
	maxTemperatureFahrenheit = prometheus.NewDesc("weather_maxtemperature_fahrenheit", "Current maximum temperature in fahrenheit", labels, nil)
	maxTemperatureKelvin     = prometheus.NewDesc("weather_maxtemperature_kelvin", "Current maximum temperature in kelvin", labels, nil)
	minTemperatureCelsius    = prometheus.NewDesc("weather_mintemperature_celsius", "Current minimum temperature in celsius", labels, nil)
	minTemperatureFahrenheit = prometheus.NewDesc("weather_mintemperature_fahrenheit", "Current minimum temperature in fahrenheit", labels, nil)
	minTemperatureKelvin     = prometheus.NewDesc("weather_mintemperature_kelvin", "Current minimum temperature in kelvin", labels, nil)
	pressureHpa              = prometheus.NewDesc("weather_pressure_hpa", "Atmospheric pressure in hectopascal", labels, nil)
	sunrise                  = prometheus.NewDesc("weather_sunrise", "Sunrise time as a UNIX timestamp", labels, nil)
	sunset                   = prometheus.NewDesc("weather_sunset", "Sunset time as a UNIX timestamp", labels, nil)
	temperatureCelsius       = prometheus.NewDesc("weather_temperature_celsius", "Current temperature in celsius", labels, nil)
	temperatureFahrenheit    = prometheus.NewDesc("weather_temperature_fahrenheit", "Current temperature in fahrenheit", labels, nil)
	temperatureKelvin        = prometheus.NewDesc("weather_temperature_kelvin", "Current temperature in kelvin", labels, nil)
	up                       = prometheus.NewDesc("weather_up", "Whether the metrics can be collected", nil, nil)
	visibilityMeters         = prometheus.NewDesc("weather_visibility_meters", "Visibility in meters", labels, nil)
	visibilityMiles          = prometheus.NewDesc("weather_visibility_miles", "Visibility in miles", labels, nil)
	windDirection            = prometheus.NewDesc("weather_winddirection", "Wind direction in degrees", labels, nil)
	windspeedMph             = prometheus.NewDesc("weather_windspeed_mph", "Wind speed in miles per hour", labels, nil)
	windspeedMps             = prometheus.NewDesc("weather_windspeed_mps", "Wind speed in meters per second", labels, nil)
)

type Exporter struct {
	apikey  string
	cityIds string
}

func newExporter(apikey string, cityIds string) *Exporter {
	return &Exporter{
		apikey:  apikey,
		cityIds: cityIds,
	}
}

func (exporter *Exporter) Describe(channel chan<- *prometheus.Desc) {
	channel <- cloudinessPercent
	channel <- feelslikeCelsius
	channel <- feelslikeFahrenheit
	channel <- feelslikeKelvin
	channel <- humidityPercent
	channel <- maxTemperatureCelsius
	channel <- maxTemperatureFahrenheit
	channel <- maxTemperatureKelvin
	channel <- minTemperatureCelsius
	channel <- minTemperatureFahrenheit
	channel <- minTemperatureKelvin
	channel <- pressureHpa
	channel <- sunrise
	channel <- sunset
	channel <- temperatureCelsius
	channel <- temperatureFahrenheit
	channel <- temperatureKelvin
	channel <- up
	channel <- visibilityMeters
	channel <- visibilityMiles
	channel <- windDirection
	channel <- windspeedMph
	channel <- windspeedMps
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

		channel <- prometheus.MustNewConstMetric(cloudinessPercent, prometheus.GaugeValue, location.Clouds.CoverPercentage, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(feelslikeCelsius, prometheus.GaugeValue, kelvinToCelsius(location.Main.FeelsLikeKelvin), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(feelslikeFahrenheit, prometheus.GaugeValue, kelvinToFahrenheit(location.Main.FeelsLikeKelvin), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(feelslikeKelvin, prometheus.GaugeValue, location.Main.FeelsLikeKelvin, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(humidityPercent, prometheus.GaugeValue, location.Main.HumidityPercent, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(maxTemperatureCelsius, prometheus.GaugeValue, kelvinToCelsius(location.Main.TemperatureMaxKelvin), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(maxTemperatureFahrenheit, prometheus.GaugeValue, kelvinToFahrenheit(location.Main.TemperatureMaxKelvin), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(maxTemperatureKelvin, prometheus.GaugeValue, location.Main.TemperatureMaxKelvin, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(minTemperatureCelsius, prometheus.GaugeValue, kelvinToCelsius(location.Main.TemperatureMinKelvin), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(minTemperatureFahrenheit, prometheus.GaugeValue, kelvinToFahrenheit(location.Main.TemperatureMinKelvin), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(minTemperatureKelvin, prometheus.GaugeValue, location.Main.TemperatureMinKelvin, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(pressureHpa, prometheus.GaugeValue, location.Main.PressureHpa, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(sunrise, prometheus.GaugeValue, location.Sys.Sunrise, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(sunset, prometheus.GaugeValue, location.Sys.Sunset, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(temperatureCelsius, prometheus.GaugeValue, kelvinToCelsius(location.Main.TemperatureKelvin), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(temperatureFahrenheit, prometheus.GaugeValue, kelvinToFahrenheit(location.Main.TemperatureKelvin), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(temperatureKelvin, prometheus.GaugeValue, location.Main.TemperatureKelvin, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(visibilityMeters, prometheus.GaugeValue, location.Visibility, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(visibilityMiles, prometheus.GaugeValue, metersToMiles(location.Visibility), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(windDirection, prometheus.GaugeValue, location.Wind.DirectionDegrees, locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(windspeedMph, prometheus.GaugeValue, meterPerSecondToMilesPerHour(location.Wind.SpeedMps), locationID, location.Name, location.Sys.Country, latitude, longitude)
		channel <- prometheus.MustNewConstMetric(windspeedMps, prometheus.GaugeValue, location.Wind.SpeedMps, locationID, location.Name, location.Sys.Country, latitude, longitude)
	}
}

func (exporter *Exporter) getData() response {
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/group?id=%s&appid=%s", exporter.cityIds, exporter.apikey))
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

	return data
}

func kelvinToCelsius(temp float64) float64 {
	return temp - 273.15
}

func kelvinToFahrenheit(temp float64) float64 {
	return kelvinToCelsius(temp)*9/5 + 32
}

func meterPerSecondToMilesPerHour(speed float64) float64 {
	return speed * 2.237
}

func metersToMiles(distance float64) float64 {
	return distance * 0.00062137
}
