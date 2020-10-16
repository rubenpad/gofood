package routes

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  string      `json:"status,omitempty"`
}

type decodeddata struct {
	History     []interface{} `json:"history,omitempty"`
	IPList      []interface{} `json:"iplist,omitempty"`
	Suggestions []interface{} `json:"suggestions,omitempty"`
}

type decodeBuyers struct {
	Buyers []interface{} `json:"buyers,omitempty"`
}
