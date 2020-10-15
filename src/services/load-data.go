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

func (ld *loadDataService) GetData() error {
	store := store.New()
	// TODO pass date as parameter to this function
	date := "1602530864"

	transactionsResponse, err := ld.makeRequest("/transactions?date=" + date)
	productsResponse, _ := ld.makeRequest("/products?date=" + date)
	buyersResponse, _ := ld.makeRequest("/buyers?date=" + date)
	if err != nil {
		log.Panic("Couldn't get data")
	}

	transactions := formatTransactionsData(transactionsResponse.Body)
	products := formatProductsData(productsResponse.Body)
	buyers := formatBuyersData(buyersResponse.Body)
	mutation := formatQueryData(transactions, products, buyers)
	encoded, _ := json.Marshal(mutation)

	savedErr := store.Save(encoded)
	if savedErr != nil {
		return savedErr
	}

	return nil
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
