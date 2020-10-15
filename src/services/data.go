package services

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/rubbenpad/gofood/store"
)

type loadDataService struct {
	httpclient *http.Client
}

func NewloadDataService() *loadDataService {
	return &loadDataService{httpclient: &http.Client{}}
}

func (ld *loadDataService) GetData(date string) (bool, error) {
	store := store.New()
	dateExists := store.GetDate(date)

	if dateExists {
		return dateExists, nil
	}

	transactionsResponse, err := ld.makeRequest("/transactions?date=" + date)
	productsResponse, _ := ld.makeRequest("/products?date=" + date)
	buyersResponse, _ := ld.makeRequest("/buyers?date=" + date)
	if err != nil {
		log.Panic("Couldn't get data")
	}

	transactions := formatTransactionsData(date, transactionsResponse.Body)
	products := formatProductsData(productsResponse.Body)
	buyers := formatBuyersData(buyersResponse.Body)
	mutation := formatQueryData(transactions, products, buyers)
	encoded, _ := json.Marshal(mutation)

	savedErr := store.Save(encoded)
	if savedErr != nil {
		return false, savedErr
	}

	return false, nil
}

func (ld *loadDataService) makeRequest(path string) (*http.Response, error) {
	url, _ := os.LookupEnv("BASE_URL")
	res, err := ld.httpclient.Get(url + path)
	if err != nil {
		log.Panic("Failed fetching data")
		return nil, err
	}

	return res, nil
}
