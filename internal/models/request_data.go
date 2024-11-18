package models

// RequestData описывает структуру данных, принимаемых от клиента.
type RequestData struct {
	Tab    string `json:"tab"`
	Status string `json:"status"`
	Data   string `json:"data"`
}
