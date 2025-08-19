package utils

import (
	"github.com/codingsince1985/geo-golang/openstreetmap"
)

func ExtractLatAndLon(location string) (float64, float64, error) {
	// Geocoder package to extract latitude and longitude from the location string
	// it is free and open source
	latAndLon, err := openstreetmap.Geocoder().Geocode(location)
	if err != nil {
		return 0, 0, err
	}
	return latAndLon.Lat, latAndLon.Lng, nil
}