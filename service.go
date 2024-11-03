package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	fullurl = "https://api.open-meteo.com/v1/forecast?latitude=59.57&longitude=30.19&hourly=temperature_2m,uv_index&past_days=1"
)

type weatherService struct {
}

func (ws *weatherService) GetWeather() (*WeatherResp, error) {
	//URL := "https://api.open-meteo.com/v1/forecast"

	u := url.URL{
		Scheme: "https",
		Host:   "api.open-meteo.com",
		Path:   "/v1/forecast",
	}

	q := url.Values{}
	q.Add("latitude", "59.57")
	q.Add("longitude", "30.19")
	q.Add("hourly", "temperature_2m")
	q.Add("past_days", "1")

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responce WeatherResp
	err = json.Unmarshal(body, &responce)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &responce, nil
}

type WeatherResp struct {
	Hourly Hourly `json:""`
}

type Hourly struct {
	Time         []string  `json:"time"`
	Temperatures []float32 `json:"temperature_2m"`
}
