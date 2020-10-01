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
	Cover float64 `json:"all"`
}

type coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type mainStruct struct {
	FeelsLike      float64 `json:"feels_like"`
	Humidity       float64 `json:"humidity"`
	Pressure       float64 `json:"pressure"`
	Temperature    float64 `json:"temp"`
	TemperatureMax float64 `json:"temp_max"`
	TemperatureMin float64 `json:"temp_min"`
}

type sys struct {
	Country string  `json:"country"`
	Sunrise float64 `json:"sunrise"`
	Sunset  float64 `json:"sunset"`
}

type wind struct {
	Direction float64 `json:"deg"`
	Speed     float64 `json:"speed"`
}
