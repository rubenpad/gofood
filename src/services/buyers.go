package services

import (
	"github.com/rubbenpad/gofood/store"
)

type buyersService struct{}

func NewBuyersService() *buyersService {
	return &buyersService{}
}

func (bs *buyersService) FindTransactions(id string) ([]byte, error) {
	store := store.New()
	data, err := store.FindTransactions(id)
	return data, err
}
