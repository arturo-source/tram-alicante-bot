package planyourroute

import (
	"encoding/json"
	"net/url"

	"github.com/arturo-source/tramalicantebot/utils"
)

type Coord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Estacion struct {
	Nombre     string      `json:"nombre"`
	Hora       string      `json:"hora"`
	EstacionID int         `json:"estacionId"`
	Mensaje    interface{} `json:"mensaje"`
	Coord      Coord       `json:"coord"`
}

type Paso struct {
	Origen            string      `json:"origen"`
	Destino           string      `json:"destino"`
	EstacionIDOrigen  int         `json:"estacionIdOrigen"`
	EstacionIDDestino int         `json:"estacionIdDestino"`
	CoordOrigen       Coord       `json:"coord_origen"`
	CoordDestino      Coord       `json:"coord_destino"`
	TrenConDestino    string      `json:"tren_con_destino"`
	Linea             int         `json:"linea"`
	LineaColor        string      `json:"linea_color"`
	NumParadas        int         `json:"num_paradas"`
	Estaciones        []Estacion  `json:"estaciones"`
	Mensaje           interface{} `json:"mensaje"`
}

type Viaje struct {
	Origen       string        `json:"origen"`
	Destino      string        `json:"destino"`
	Fecha        string        `json:"fecha"`
	HoraInicio   string        `json:"hora_inicio"`
	HoraFin      string        `json:"hora_fin"`
	Duracion     int           `json:"duracion"`
	Zonas        string        `json:"zonas"`
	Lineas       []int         `json:"lineas"`
	Pasos        []Paso        `json:"pasos"`
	Avisos       []interface{} `json:"avisos"`
	Mensaje      string        `json:"mensaje"`
	CoordOrigen  Coord         `json:"coord_origen"`
	CoordDestino Coord         `json:"coord_destino"`
}

type PlanedTramResponse struct {
	Error   bool    `json:"error"`
	Mensaje string  `json:"mensaje"`
	Data    []Viaje `json:"data"`
	// HTML    string  `json:"html"`
}

func Route(from, to, day, hour string) (PlanedTramResponse, error) {
	data := url.Values{}
	data.Set("auth_token", "b185e97c90d663ece79e3bc794ce8163")
	data.Set("action", "planificar-ruta")
	data.Set("origen", from)
	data.Set("destino", to)
	data.Set("dia", day)
	data.Set("salida", "0") // 0 = envía hora de salida, 1 = envía hora de llegada
	data.Set("hora", hour)

	v := url.Values{}
	v.Set("action", "formularios_ajax")
	v.Set("data", data.Encode())

	body, err := utils.DoRequest(v)
	if err != nil {
		return PlanedTramResponse{}, err
	}

	var dataResp PlanedTramResponse
	err = json.Unmarshal(body, &dataResp)
	return dataResp, err
}
