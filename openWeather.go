package main

import (
	"fmt"
	"strings"
	"errors"
	"encoding/json"
	"time"
)

type OpenWeatherCondition struct {
	Main string
	Description string
}

type OpenWeatherResponseCurrent struct {
	Dt int64
	Temp float32
	Feels_like float32
	Pressure int
	Humidity int
	Weather []OpenWeatherCondition

}

type OpenWeatherResponseHourly struct {
	Dt int64
	Temp float32
	Feels_like float32
	Pressure int
	Humidity int
	Weather []OpenWeatherCondition
}

type OpenWeatherResponseDaily struct {
	Dt int64
	Sunrise int64
	Sunset int64
	Temp struct {
		Day float32
		Min float32
		Max float32
		Night float32
		Eve float32
		Morn float32
	}
	Feels_like struct {
		Day float32
		Night float32
		Eve float32
		Morn float32
	}
	Pressure int
	Humidity int
	Weather []OpenWeatherCondition
}


type OpenWeatherResponse struct {
	Current *OpenWeatherResponseCurrent
	Hourly *[]OpenWeatherResponseHourly
	Daily *[]OpenWeatherResponseDaily
}

func (weather OpenWeatherResponseCurrent) Output(units string) string {
	var unitAbbr string

	switch units {
		case "metric":
			unitAbbr = " C"
		case "imperial":
			unitAbbr = " F"
	}

	return fmt.Sprintf("Current: %g%s | Feels Like: %5.2f%s | Humidity: %d%% | %s\n",
		weather.Temp,
		unitAbbr,
		weather.Feels_like,
		unitAbbr,
		weather.Humidity,
		weather.Weather[0].Description,
	)
}

func (weather OpenWeatherResponseHourly) Output(units string) string {
	var unitAbbr string

	switch units {
		case "metric":
			unitAbbr = "C"
		case "imperial":
			unitAbbr = "F"
	}

	t := time.Unix(weather.Dt, 0)
	return fmt.Sprintf("%-9s %2d/%2d %02d:00    %5.2f%s | Feels Like: %5.2f%s |  Humidity: %d%% | %s\n",
		t.Weekday().String(),
		t.Month(),
		t.Day(),
		t.Hour(),
		weather.Temp,
		unitAbbr,
		weather.Feels_like,
		unitAbbr,
		weather.Humidity,
		weather.Weather[0].Description,
	)
}

func (weather OpenWeatherResponseDaily) Output(units string) string {
	var unitAbbr string

	switch units {
		case "metric":
			unitAbbr = "C"
		case "imperial":
			unitAbbr = "F"
	}

	t := time.Unix(weather.Dt, 0)
	return fmt.Sprintf("%-9s %2d/%2d  High:%5.2f%s Low: %5.2f%s |  Humidity: %d%% | %s\n",
		t.Weekday().String(),
		t.Month(),
		t.Day(),
		weather.Temp.Max,
		unitAbbr,
		weather.Temp.Min,
		unitAbbr,
		weather.Humidity,
		weather.Weather[0].Description,
	)
}

func weatherOfLocationInPeriod(location LatLngStruct, period string, unit string) (weather OpenWeatherResponse, responseError error) {
	var possiblePeriods = [4]string{"current","minutely","hourly","daily"}
	exclude := make([]string, 0)

	for i := 0; i < len(possiblePeriods); i++ {
		if possiblePeriods[i] != period {
			exclude = append(exclude, possiblePeriods[i])
		}
	}

	excludedString := strings.Join(exclude, ",")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%g&lon=%g&exclude=%s&appid=%s&units=%s",
		location.Lat,
		location.Lng,
		excludedString,
		OpenWeatherApiKey,
		unit,
	)

	response, requestError := httpClient.Get(url)
	if requestError != nil {
		return weather, requestError
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return weather, errors.New(fmt.Sprintf("OpenWeatherRequestFailed: %s", response.Status))
	}

	requestError = json.NewDecoder(response.Body).Decode(&weather)

	return weather, requestError
}