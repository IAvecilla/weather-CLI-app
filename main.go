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

	location, requestError := latittudeAndLongitudOfPlace("denver, co")
	if requestError != nil {
		panic(requestError)
	}

	fmt.Printf("%+v\n", location)
}