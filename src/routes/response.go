package routes

import (
	"github.com/rubbenpad/gofood/domain"
)

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Status  string      `json:"status,omitempty"`
}

type decodeddata struct {
	Buyer       interface{}   `json:"buyer"`
	History     []interface{} `json:"history"`
	IPList      []interface{} `json:"iplist"`
	Suggestions []interface{} `json:"suggestions"`
}

type decodeBuyers struct {
	Buyers []domain.Buyer `json:"buyers"`
}
