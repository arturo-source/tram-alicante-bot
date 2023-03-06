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
	Lines []struct {
		LineName string
		Hours    []string
	}
}

func newLine(lineName string, hours []string) struct {
	LineName string
	Hours    []string
} {
	return struct {
		LineName string
		Hours    []string
	}{
		LineName: lineName,
		Hours:    hours,
	}
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

	extractLineName := func(line string) string {
		startLineNameToken := `<span class="c-negro fs--18 fw--600">`
		startLineName := strings.Index(line, startLineNameToken) + len(startLineNameToken)
		endLineNameToken := `</span>`
		endLineName := strings.Index(line[startLineName:], endLineNameToken)

		return line[startLineName : startLineName+endLineName]
	}
	extractHours := func(line string) []string {
		var hours []string

		hoursHtmlToken := `<div class="df-s">`
		hoursHtmlSplited := strings.Split(line, hoursHtmlToken)
		for i := 2; i < len(hoursHtmlSplited); i++ {
			startHourToken := `<p class="hora">`
			startHour := strings.Index(hoursHtmlSplited[i], startHourToken) + len(startHourToken)
			endHourToken := `</p>`
			endHour := strings.Index(hoursHtmlSplited[i][startHour:], endHourToken)
			hour := hoursHtmlSplited[i][startHour : startHour+endHour]

			minuteToken := `<p class="minuto`
			minuteSplited := strings.Split(hoursHtmlSplited[i], minuteToken)
			for j := 0; j < len(minuteSplited); j++ {
				startMinuteToken := `">`
				startMinute := strings.Index(minuteSplited[j], startMinuteToken) + len(startMinuteToken)
				endMinuteToken := `</p>`
				endMinute := strings.Index(minuteSplited[j][startMinute:], endMinuteToken)
				minute := minuteSplited[j][startMinute : startMinute+endMinute]

				hours = append(hours, hour+":"+minute)
			}

		}

		return hours
	}

	lines := strings.Split(html, `<div class="row separador-rojo-o toggle--tren-destino collapsed"`)
	for i := 1; i < len(lines); i++ {
		lineName := extractLineName(lines[i])
		hours := extractHours(lines[i])

		ss.Lines = append(ss.Lines, newLine(lineName, hours))
	}

	return ss
}
