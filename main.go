package main
import (
	"net/http"
	"time"
	"fmt"
)

var httpClient http.Client

func main() {
	httpClient = http.Client{
		Timeout: time.Second * 10,
	}

	location, requestError := latittudeAndLongitudOfPlace("New York")
	if requestError != nil {
		panic(requestError)
	}

	weather, requestError := weatherOfLocationInPeriod(location, "hourly", "metric")
	if requestError != nil {
		panic(requestError)
	}

	printWeatherResult(*weather.Hourly, "New York, EEUU", "metric")
}

func printWeatherResult (weather interface{}, place string, units string) {
	fmt.Printf("Weather for %s:\n", place)

	switch weather.(type) {
	case OpenWeatherResponseCurrent:
		fmt.Printf(weather.(OpenWeatherResponseCurrent).Output(units)) 

	case []OpenWeatherResponseHourly:
		for _, h := range weather.([]OpenWeatherResponseHourly) {
			fmt.Print(h.Output(units))
		}

	case []OpenWeatherResponseDaily:
		for _, h := range weather.([]OpenWeatherResponseDaily) {
			fmt.Print(h.Output(units))
		}
	}
}