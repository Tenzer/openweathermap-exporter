package main

type response struct {
	Locations []location `json:"list"`
}

type location struct {
	Clouds      clouds      `json:"clouds"`
	Coordinates coordinates `json:"coord"`
	ID          int64       `json:"id"`
	Main        mainStruct  `json:"main"`
	Name        string      `json:"name"`
	Sys         sys         `json:"sys"`
	Visibility  float64     `json:"visibility"`
	Wind        wind        `json:"wind"`
}

type clouds struct {
	CoverPercentage float64 `json:"all"`
}

type coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type mainStruct struct {
	FeelsLikeKelvin      float64 `json:"feels_like"`
	HumidityPercent      float64 `json:"humidity"`
	PressureHpa          float64 `json:"pressure"`
	TemperatureKelvin    float64 `json:"temp"`
	TemperatureMaxKelvin float64 `json:"temp_max"`
	TemperatureMinKelvin float64 `json:"temp_min"`
}

type sys struct {
	Country string  `json:"country"`
	Sunrise float64 `json:"sunrise"`
	Sunset  float64 `json:"sunset"`
}

type wind struct {
	DirectionDegrees float64 `json:"deg"`
	SpeedMps         float64 `json:"speed"`
}
