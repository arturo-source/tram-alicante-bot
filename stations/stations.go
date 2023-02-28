package stations

import (
	"encoding/json"
	"os"
	"strings"
)

type Estacion struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func AllStations() ([]Estacion, error) {
	const stopsFileName = "stations/stations.json"
	stops := []Estacion{}

	file, err := os.ReadFile(stopsFileName)
	if err != nil {
		return stops, err
	}
	err = json.Unmarshal(file, &stops)
	if err != nil {
		return stops, err
	}

	return stops, nil
}

func SimilarStations(station string) ([]Estacion, error) {
	stops, err := AllStations()
	if err != nil {
		return nil, err
	}
	station = strings.ToLower(station)
	stationWords := strings.Split(station, " ")

	isSearchable := func(word string) bool {
		return len(word) > 3
	}

	var similarStops []Estacion

	for _, stop := range stops {
		stopName := strings.ToLower(stop.Name)

		for _, word := range stationWords {
			if !isSearchable(word) {
				continue
			}

			if strings.Contains(stopName, word) {
				similarStops = append(similarStops, stop)
			}
		}
	}

	return similarStops, nil
}
