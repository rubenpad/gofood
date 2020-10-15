package routes

type response struct {
	Data    []interface{} `json:"data,omitempty"`
	Message string        `json:"message,omitempty"`
	Status  string        `json:"status,omitempty"`
}
