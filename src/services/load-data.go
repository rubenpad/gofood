package services

import (
	"log"
	"net/http"
)

// TODO set environment variables
const baseURL = "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com"

type loadDataService struct {
	url        string
	httpclient *http.Client
}

func NewloadDataService() *loadDataService {
	return &loadDataService{url: baseURL, httpclient: &http.Client{}}
}

func (ld *loadDataService) GetData() queryMutation {
	// TODO pass date as parameter to this function
	date := "1602530864"

	transactionsResponse, err := ld.makeRequest("/transactions?date=" + date)
	productsResponse, _ := ld.makeRequest("/products?date=" + date)
	buyersResponse, _ := ld.makeRequest("/buyers?date=" + date)
	if err != nil {
		log.Panic("No")
	}

	transactions := formatTransactionsData(transactionsResponse.Body)
	products := formatProductsData(productsResponse.Body)
	buyers := formatBuyersData(buyersResponse.Body)
	queryset := formatQueryData(transactions, products, buyers)

	return queryset
}

func (ld *loadDataService) makeRequest(path string) (*http.Response, error) {
	res, err := ld.httpclient.Get(ld.url + path)
	if err != nil {
		log.Fatal("Failed fetching data")
		return nil, err
	}

	return res, nil
}
