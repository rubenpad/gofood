package services

import (
	"github.com/rubbenpad/gofood/store"
)

type buyersService struct{}

func NewBuyersService() *buyersService {
	return &buyersService{}
}

func (bs *buyersService) FindTransactions(id string) {
	store := store.New()
	store.FindTransactions(id)
	//return data, err
}
