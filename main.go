package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	q := "-20.81972,-49.37944"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=ec180872243c4f57a4f153631230105&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		fmt.Println("Error fetching weather data:", err)
		return
	}
	defer res.Body.Close()

	fmt.Println("Response Status:", res.Status)
	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s: %.0f°C, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	fmt.Println("Hourly Forecast:")
	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Hour() == time.Now().Hour() {
			fmt.Printf("%s - %.0f°C, %.0f%%, %s\n",
				date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)
		}
	}
}
