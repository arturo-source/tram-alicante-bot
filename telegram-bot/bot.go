package telegrambot

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	planyourroute "github.com/arturo-source/tramalicantebot/plan-your-route"
	routeschedules "github.com/arturo-source/tramalicantebot/route-schedules"
	stationschedules "github.com/arturo-source/tramalicantebot/station-schedules"
	stations "github.com/arturo-source/tramalicantebot/stations"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	STATION_SCHEDULES = "Horarios de la estación"
	ROUTE_SCHEDULES   = "Horarios de la ruta"
	PLAN_ROUTE        = "Planificar ruta"

	END_LINE = "\n"
)

func estaciones() string {
	var text string
	stations, _ := stations.AllStations()
	for _, s := range stations {
		text += s.Name + END_LINE
	}
	return text
}

func horarios(args string, day string) []string {
	var texts []string
	stationName := args
	stations, err := stations.MostSimilarStations(stationName)
	if err != nil {
		texts = append(texts, "No se pudo obtener las estaciones")
		return texts
	}

	if len(stations) == 0 {
		text := "No se encontró ninguna estación con ese nombre"
		texts = append(texts, text)
	}

	if len(stations) > 1 {
		text := "Varias estaciones coinciden con ese nombre:\n"
		for _, s := range stations {
			text += s.Name + END_LINE
		}
		texts = append(texts, text)
	}

	for _, station := range stations {
		text := STATION_SCHEDULES + END_LINE
		ss, _ := stationschedules.Schedules(station.Id, day)

		for _, line := range ss.Lines {
			text += line.LineName + END_LINE
			for _, hour := range line.Hours {
				text += "- " + hour + END_LINE
			}
			text += END_LINE
		}

		texts = append(texts, text)
	}

	return texts
}

func getStationName(station string) (string, error) {
	stationsFrom, _ := stations.MostSimilarStations(station)
	if len(stationsFrom) == 0 {
		return "", fmt.Errorf("No se encontró ninguna estación con el nombre %s", station)
	}

	if len(stationsFrom) > 1 {
		textTemplate := `No se puede obtener la ruta porque varias estaciones coinciden con el nombre "%s":` + END_LINE
		for _, s := range stationsFrom {
			textTemplate += s.Name + END_LINE
		}

		return "", fmt.Errorf(textTemplate, station)
	}

	return stationsFrom[0].Id, nil
}

func ruta(args, day, hour string) []string {
	var texts []string
	route := strings.Split(args, "-")

	if len(route) < 2 {
		text := `Debes especificar la ruta con el formato "/ruta origen - destino"`
		texts = append(texts, text)
		return texts
	}

	_, err := stations.MostSimilarStations("")
	if err != nil {
		texts = append(texts, "No se pudo obtener las estaciones")
		return texts
	}

	from := strings.TrimSpace(route[0])
	to := strings.TrimSpace(route[1])
	fromId, err := getStationName(from)
	if err != nil {
		texts = append(texts, err.Error())
	}
	toId, err := getStationName(to)
	if err != nil {
		texts = append(texts, err.Error())
	}

	if fromId == "" || toId == "" {
		return texts
	}

	s, _ := routeschedules.Route(fromId, toId, day, hour)
	if s.Mensaje != "" {
		texts = append(texts, "Hay un mensaje: "+s.Mensaje)
	}
	if s.Avisos != "" {
		texts = append(texts, "Hay un aviso: "+s.Avisos)
	}

	textTemplate := ROUTE_SCHEDULES + END_LINE + `Desde %s hasta %s, se tarda ~%d minutos y son ~%d metros. Recorrerás las zonas %s.`
	if len(s.Result) > 1 {
		textTemplate += fmt.Sprintf("Encontrarás %d opciones.", len(s.Result))
	}
	texts = append(texts, fmt.Sprintf(textTemplate, s.Origen.Nombre, s.Destino.Nombre, s.Duracion/60, s.Distancia, s.Zonas))

	if len(s.Result) > 1 {
		for i, option := range s.Result {
			if len(option.Pasos) == 1 {
				textTemplate = `Opción %d (sin transbordos):` + END_LINE + `Puedes coger trenes con destino a %s.`
				texts = append(texts, fmt.Sprintf(textTemplate, i+1, option.Pasos[0].TrenesConDestino))
			} else if len(option.Pasos) > 1 {
				textTemplate = `Opción %d:` + END_LINE
				for _, paso := range option.Pasos {
					textTemplate += "Sube en " + paso.Origen + " y baja en " + paso.Destino + ". Sube a un tren con algún destino " + paso.TrenesConDestino + ".\n"
				}
				texts = append(texts, fmt.Sprintf(textTemplate, i+1))
			}
		}
	}

	for _, horario := range s.Horarios {
		if len(horario.Horas) == 0 {
			textTemplate = "Los trenes con destinos %s no saldrán más el %s a partir de las %s."
			texts = append(texts, fmt.Sprintf(textTemplate, horario.Destinos, day, hour))
			continue
		}

		var hours []string
		for _, hour := range horario.Horas {
			hours = append(hours, hour[0])
		}
		textTemplate = "Los trenes con destinos %s salen a las siguientes horas:\n%s"
		texts = append(texts, fmt.Sprintf(textTemplate, horario.Destinos, strings.Join(hours, END_LINE)))
	}

	return texts
}

func siguiente(args, day, hour string) []string {
	var texts []string
	route := strings.Split(args, "-")

	if len(route) < 2 {
		text := `Debes especificar la ruta con el formato "/siguiente origen - destino"`
		texts = append(texts, text)
		return texts
	}

	_, err := stations.MostSimilarStations("")
	if err != nil {
		texts = append(texts, "No se pudo obtener las estaciones")
		return texts
	}

	from := strings.TrimSpace(route[0])
	to := strings.TrimSpace(route[1])
	fromId, err := getStationName(from)
	if err != nil {
		texts = append(texts, err.Error())
	}
	toId, err := getStationName(to)
	if err != nil {
		texts = append(texts, err.Error())
	}

	if fromId == "" || toId == "" {
		return texts
	}

	r, err := planyourroute.Route(fromId, toId, day, hour)
	if err != nil {
		texts = append(texts, "No se pudo obtener la ruta")
		return texts
	}

	if len(r.Data) == 0 {
		texts = append(texts, "No hay trenes disponibles.")
		return texts
	}

	textTemplate := PLAN_ROUTE + END_LINE + `El siguiente tram desde %s hasta %s, tiene una duración aproximada de %d min. Recorrerás las zonas %s. `
	if len(r.Data) > 1 {
		textTemplate += fmt.Sprintf("Encontrarás %d opciones.", len(r.Data))
	}
	firstOpt := r.Data[0]
	texts = append(texts, fmt.Sprintf(textTemplate, firstOpt.Origen, firstOpt.Destino, firstOpt.Duracion, firstOpt.Zonas))

	for i, option := range r.Data {
		text := fmt.Sprintf(`Opción %d: %d transbordos. Sale a las %s, y llega a las %s.`+END_LINE, i+1, len(option.Pasos)-1, option.HoraInicio, option.HoraFin)
		for j, paso := range option.Pasos {
			text += fmt.Sprintf("1. Sube al tren con destino %s en la parada %s."+END_LINE, paso.TrenConDestino, paso.Origen)

			text += "2. Pasarás por las siguientes paradas: "
			for _, estacion := range paso.Estaciones {
				text += fmt.Sprint(estacion.Nombre, " (", estacion.Hora, "), ")
			}
			text += END_LINE

			if j == len(option.Pasos)-1 {
				text += fmt.Sprintf("3. Baja en la parada %s.", paso.Destino)
			} else {
				text += fmt.Sprintf("3. Baja en la parada %s.", option.Pasos[j+1].Origen)
			}

			text += END_LINE + END_LINE
		}
		texts = append(texts, text)
	}

	return texts
}

func responseCommand(msg *tgbotapi.Message) []string {
	var texts []string
	currDay := time.Now().Format("2006-01-02")
	currHour := time.Now().Format("15:04")
	args := msg.CommandArguments()

	switch msg.Command() {
	case "estaciones":
		resp := estaciones()
		texts = append(texts, resp)
	case "horarios":
		if args == "" {
			text := `El formato correcto es "/horarios nombre_estacion", por ejemplo "/horarios luceros"`
			texts = append(texts, text)
			break
		}

		resp := horarios(args, currDay)
		texts = append(texts, resp...)
	case "ruta":
		if args == "" {
			text := `El formato correcto es "/ruta origen - destino", por ejemplo "/ruta luceros - universidad"`
			texts = append(texts, text)
			break
		}

		resp := ruta(args, currDay, currHour)
		texts = append(texts, resp...)
	case "siguiente":
		if args == "" {
			text := `El formato correcto es "/siguiente origen - destino", por ejemplo "/siguiente luceros - universidad"`
			texts = append(texts, text)
			break
		}

		resp := siguiente(args, currDay, currHour)
		texts = append(texts, resp...)
	default:
		text :=
			`Comandos disponibles:
/estaciones -> Lista de nombres de las estaciones.
/horarios luceros -> Lista los horarios de luceros hoy.
/ruta luceros - universidad -> Lista los horarios de esta ruta.
/siguiente luceros - universidad -> Muestra la hora del siguiente tram que va de luceros a la universidad.

Además puedes responder a los mensajes enviados por mí poniendo la fecha "%s", la hora "%s", o la fecha y hora "%s %s". Por defecto se usará la fecha y hora actual.
`
		texts = append(texts, fmt.Sprintf(text, currDay, currHour, currDay, currHour))
	}

	return texts
}

// func repliedMessage(msg *tgbotapi.Message) string {
// 	var text string

// 	firstLine := strings.Split(msg.ReplyToMessage.Text, "\n")[0]

// 	switch firstLine {
// 	case STATION_SCHEDULES:
// 		text = "Horarios de la estación"
// 	case ROUTE_SCHEDULES:
// 		text = "Horarios de la ruta"
// 	case PLAN_ROUTE:
// 		text = "Planificar ruta"
// 	default:
// 		text = "I don't know that command"
// 	}

// 	return text
// }

func Run() error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		var texts []string

		// if update.Message.ReplyToMessage != nil {
		// 	text := repliedMessage(update.Message)
		// 	texts = append(texts, text)
		// }

		if update.Message.IsCommand() {
			texts = responseCommand(update.Message)
		}

		if len(texts) == 0 {
			continue
		}

		for _, text := range texts {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			if _, err := bot.Send(msg); err != nil {
				return err
			}
		}
	}

	return nil
}
