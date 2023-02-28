package routeschedules

import (
	"encoding/json"
	"net/url"

	"github.com/arturo-source/tramalicantebot/utils"
)

type Destino struct {
	Destino string      `json:"destino"`
	Linea   int         `json:"linea"`
	ID      interface{} `json:"id"`
}

type Parada struct {
	EstacionID int    `json:"estacionId"`
	Nombre     string `json:"nombre"`
}

type SchedulesTramResponse struct {
	Origen        Parada  `json:"origen"`
	Destino       Parada  `json:"destino"`
	Fecha         string  `json:"fecha"`
	Salida        string  `json:"salida"`
	Llegada       string  `json:"llegada"`
	Zonas         string  `json:"zonas"`
	Duracion      int     `json:"duracion"`
	Distancia     int     `json:"distancia"`
	HuellaCarbono float64 `json:"huellaCarbono"`
	Accesibilidad int     `json:"accesibilidad"`
	Transbordos   []struct {
		Inicio     json.Number `json:"inicio"`
		EstInicio  string      `json:"estInicio"`
		Fin        json.Number `json:"fin"`
		EstFin     string      `json:"estFin"`
		Tren       string      `json:"tren"`
		TrenPos    int         `json:"trenPos"`
		TripID     int         `json:"tripId"`
		Salida     int         `json:"salida"`
		Llegada    int         `json:"llegada"`
		ViaOrigen  string      `json:"viaOrigen"`
		ViaDestino string      `json:"viaDestino"`
		Hsalida    string      `json:"Hsalida"`
		Hllegada   string      `json:"Hllegada"`
		HoraActual string      `json:"horaActual"`
	} `json:"transbordos"`
	TransbordosPreferentes int `json:"transbordosPreferentes"`
	Horarios               []struct {
		Inicio   Parada          `json:"inicio"`
		Fin      Parada          `json:"fin"`
		Horas    [][]string      `json:"horas"`
		Trenes   map[int]Destino `json:"trenes,omitempty"`
		Destinos string          `json:"destinos"`
		Mensaje  interface{}     `json:"mensaje"`
	} `json:"horarios"`
	Avisos  string `json:"avisos"`
	Mensaje string `json:"mensaje"`
	Result  []struct {
		Duracion   int           `json:"duracion"`
		Origen     string        `json:"origen"`
		Destino    string        `json:"destino"`
		Dia        string        `json:"dia"`
		HoraInicio string        `json:"hora_inicio"`
		HoraFin    string        `json:"hora_fin"`
		Zonas      string        `json:"zonas"`
		Lineas     []interface{} `json:"lineas"`
		Pasos      []struct {
			Origen           string              `json:"origen"`
			Destino          string              `json:"destino"`
			TrenesConDestino string              `json:"trenes_con_destino"`
			Horarios         map[string][]string `json:"horarios"`
		} `json:"pasos"`
		Mensaje string `json:"mensaje"`
	} `json:"result"`
	// HTML string `json:"html"`
}

func Route(from, to, day, hour string) (SchedulesTramResponse, error) {
	data := url.Values{}
	data.Set("auth_token", "b185e97c90d663ece79e3bc794ce8163")
	data.Set("action", "horarios-ruta")
	data.Set("origen", from)
	data.Set("destino", to)
	data.Set("dia", day)
	data.Set("horaDesde", hour)
	data.Set("horaHasta", "23:59")

	v := url.Values{}
	v.Set("action", "formularios_ajax")
	v.Set("data", data.Encode())

	body, err := utils.DoRequest(v)
	if err != nil {
		return SchedulesTramResponse{}, err
	}

	var dataResp SchedulesTramResponse
	err = json.Unmarshal(body, &dataResp)
	return dataResp, err
}
