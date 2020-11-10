package services

import (
	"encoding/json"

	"github.com/rubbenpad/gofood/domain"
	"github.com/rubbenpad/gofood/store"
)

type dataService struct{}

func NewDataService() *dataService {
	return &dataService{}
}

func (ld *dataService) Load(date string) (bool, error) {
	store := store.New()
	etl := NewETLService()

	if dateExists := store.GetDate(date); dateExists {
		return dateExists, nil
	}

	// Save date
	d := domain.Timestamp{UID: "_:" + date, Date: date}
	encodedDate, _ := json.Marshal(d)
	assignedDate, _ := store.Save(encodedDate)

	go etl.GetData(assignedDate.Uids[date], date)

	return false, nil
}
