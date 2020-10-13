package services

import (
	"log"
	"net/http"

	"github.com/rubbenpad/gofood/domain"
)

const baseURL = "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com"

type loadDataService struct {
	url        string
	httpclient *http.Client
}

func NewloadDataService() *loadDataService {
	return &loadDataService{url: baseURL, httpclient: &http.Client{}}
}

func (ld *loadDataService) GetData() []domain.Transaction {
	transactionsResponse, err := ld.makeRequest("/transactions?date=" + "1602530864")
	if err != nil {
		log.Fatal(err)
	}

	transactions := formatNoStandardData(transactionsResponse.Body)
	return transactions
}

func (ld *loadDataService) makeRequest(path string) (*http.Response, error) {
	res, err := ld.httpclient.Get(ld.url + path)
	if err != nil {
		log.Fatal("Failed fetching data")
		return nil, err
	}

	return res, nil
}
