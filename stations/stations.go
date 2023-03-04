package stations

import (
	"encoding/json"
	"os"
	"strings"
)

type Estacion struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Alternative string `json:"alternative"`

	levDistance int
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

func MostSimilarStations(station string) ([]Estacion, error) {
	var mostSimilars []Estacion
	var minDistance int
	stops, err := AllStations()
	if err != nil {
		return mostSimilars, err
	}

	for i := range stops {
		levName := levenshteinDistance(station, stops[i].Name)
		levAlter := levenshteinDistance(station, stops[i].Alternative)

		stops[i].levDistance = min(levName, levAlter)

		if i == 0 || stops[i].levDistance < minDistance {
			minDistance = stops[i].levDistance
		}
	}

	for _, s := range stops {
		if s.levDistance == minDistance {
			mostSimilars = append(mostSimilars, s)
		}
	}

	return mostSimilars, nil
}

func levenshteinDistance(s, t string) int {
	s = strings.ToLower(s)
	t = strings.ToLower(t)

	m := len(s)
	n := len(t)

	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		d[i][0] = i
	}

	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			if s[i-1] == t[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+1)
			}
		}
	}

	return d[m][n]
}

func min(n ...int) int {
	min := n[0]
	for _, i := range n {
		if i < min {
			min = i
		}
	}
	return min
}
