package main

import (
	"fmt"
	"encoding/json"
	"net/url"
)

type LatLngStruct struct {
	Lat float64
	Lng float64
}

type GoogleGeocodeResult struct {
	Geometry struct {
		Location LatLngStruct
	}
}

type GoogleGeocodeResponse struct {
	Status  string
	Results []GoogleGeocodeResult
}

func latittudeAndLongitudOfPlace(place string) (latLng LatLngStruct, requestError error) {
	escPlace := url.QueryEscape(place)
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?key=%s&address=%s",
		GoogleApiKey,
		escPlace,
	)

	response, requestError := httpClient.Get(url)
	if requestError != nil {
		return latLng, requestError
	}

	defer response.Body.Close()
	
	var geocode GoogleGeocodeResponse

	requestError = json.NewDecoder(response.Body).Decode(&geocode)
	if requestError != nil {
		return latLng, requestError
	}

	if geocode.Status != "OK" || len(geocode.Results) < 1 {
		return latLng, requestError
	}

	return geocode.Results[0].Geometry.Location, requestError
}

