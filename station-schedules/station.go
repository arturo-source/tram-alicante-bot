package stationschedules

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/arturo-source/tramalicantebot/utils"
)

type response struct {
	HTML string `json:"html"`
}

type StationSchedule struct {
	Lines []Line
}

type Line struct {
	LineName string
	Hours    []string
}

func Schedules(stationId, day string) (StationSchedule, error) {
	data := url.Values{}
	data.Set("auth_token", "b185e97c90d663ece79e3bc794ce8163")
	data.Set("action", "horarios-estacion")
	data.Set("estacion", stationId)
	data.Set("dia", day)

	v := url.Values{}
	v.Set("action", "formularios_ajax")
	v.Set("data", data.Encode())

	body, err := utils.DoRequest(v)
	if err != nil {
		return StationSchedule{}, err
	}

	var dataResp response
	err = json.Unmarshal(body, &dataResp)
	if err != nil {
		return StationSchedule{}, err
	}

	if strings.Contains(dataResp.HTML, "No hi ha resultats disponibles") {
		currDay := time.Now()
		err = fmt.Errorf(
			`Fecha mal formateada o no hay horarios para esa fecha. El formato estándar es año-mes-día, por ejemplo: %s pero existen otros:
aaaa/mm/dd -> %s
aaaa-m-d ---> %s
mm-dd ------> %s
`, currDay.Format("2006-01-02"), currDay.Format("2006/01/02"), currDay.Format("2006-1-2"), currDay.Format("01-02"))
		return StationSchedule{}, err
	}

	ss := scrapeStationSchedules(dataResp.HTML)
	return ss, nil
}

func scrapeStationSchedules(html string) StationSchedule {
	var ss StationSchedule

	extractBetweenTokens := func(content, startToken, endToken string) string {
		startIndex := strings.Index(content, startToken) + len(startToken)
		endIndex := strings.Index(content[startIndex:], endToken)

		return content[startIndex : startIndex+endIndex]
	}

	extractLineName := func(line string) string {
		startLineNameToken := `<span class="c-negro fs--18 fw--600">`
		endLineNameToken := `</span>`

		return extractBetweenTokens(line, startLineNameToken, endLineNameToken)
	}

	extractHours := func(line string) []string {
		var hours []string

		hoursHtmlToken := `<div class="df-s">`
		hoursSplited := strings.Split(line, hoursHtmlToken)

		for i := 2; i < len(hoursSplited); i++ {
			startHourToken := `<p class="hora">`
			endHourToken := `</p>`
			hour := extractBetweenTokens(hoursSplited[i], startHourToken, endHourToken)

			minutesToken := `<p class="minuto`
			minutesSplited := strings.Split(hoursSplited[i], minutesToken)

			for j := 1; j < len(minutesSplited); j++ {
				startMinuteToken := `">`
				endMinuteToken := `</p>`
				minute := extractBetweenTokens(minutesSplited[j], startMinuteToken, endMinuteToken)

				hours = append(hours, hour+":"+minute)
			}

		}

		return hours
	}

	lines := strings.Split(html, `<div class="row separador-rojo-o toggle--tren-destino collapsed"`)
	for i := 1; i < len(lines); i++ {
		lineName := extractLineName(lines[i])
		hours := extractHours(lines[i])

		ss.Lines = append(ss.Lines, Line{lineName, hours})
	}

	return ss
}
