package services

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const baseURL = "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com"

type LoadDataService struct {
	url        string
	httpclient *http.Client
}

func NewLoadDataService() *LoadDataService {
	return &LoadDataService{url: baseURL, httpclient: &http.Client{}}
}

func (ld *LoadDataService) GetData() {
	res, err := ld.makeRequest("/buyers?date=" + "1602530864")
	if err != nil {
		log.Fatal(err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%s", data)
}

func (ld *LoadDataService) makeRequest(path string) (*http.Response, error) {
	res, err := ld.httpclient.Get(ld.url + path)
	if err != nil {
		log.Fatal("Failed fetching data")
		return nil, err
	}

	return res, nil
}

func (ld *LoadDataService) formatNoStandarData(content io.ReadCloser) {

}
