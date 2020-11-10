package services

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rubbenpad/gofood/store"
)

type etlService struct {
	httpclient *http.Client
}

func NewETLService() *etlService {
	return &etlService{httpclient: &http.Client{}}
}

func (etl *etlService) GetData(dateUID, date string) {
	store := store.New()

	// Build requests to remote data and fetch concurrently
	requests := etl.buildRequests(date)
	results := etl.fetchConcurrently(requests)

	// Format, encode and save products and buyers data
	all := store.FindAll()
	products := formatProductsData(results["products"].response.data, all.Products)
	buyers := formatBuyersData(results["buyers"].response.data, all.Buyers)

	assignedProducts, _ := store.Save(products)
	assignedBuyers, _ := store.Save(buyers)

	// Format, encode and save transactions data
	transactions := formatTransactionsData(
		dateUID,
		results["transactions"].response.data,
		assignedProducts.Uids,
		assignedBuyers.Uids,
	)

	if _, err := store.Save(transactions); err != nil {
		log.Panic("Error trying to store transactions data")
	}
}

func (etl *etlService) buildRequests(date string) map[string]func() (*remoteResponse, error) {
	baseurl, _ := os.LookupEnv("BASE_URL")
	endpoints := map[string]string{
		"transactions": "/transactions?date=",
		"products":     "/products?date=",
		"buyers":       "/buyers?date=",
	}

	requests := make(map[string]func() (*remoteResponse, error))
	for i := range endpoints {
		endpoint, key := endpoints[i], i
		requests[key] = func() (*remoteResponse, error) {
			return etl.makeRequest(baseurl + endpoint + date)
		}
	}

	return requests
}

type remoteResponse struct {
	data []byte
}

func (etl *etlService) makeRequest(url string) (*remoteResponse, error) {
	res, err := etl.httpclient.Get(url)
	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(res.Body)
	return &remoteResponse{data: data}, nil
}

type requestResult struct {
	response *remoteResponse
	err      error
	key      string
}

func (etl *etlService) fetchConcurrently(requests map[string]func() (*remoteResponse, error)) map[string]*requestResult {
	cn := make(chan *requestResult, len(requests))
	fns := make([]func(), len(requests))

	i := 0
	for k := range requests {
		f, key := requests[k], k
		fns[i] = func() {
			res, err := f()
			cn <- &requestResult{response: res, err: err, key: key}
		}
		i++
	}

	callConcurrent(fns)
	close(cn)

	results := make(map[string]*requestResult)
	for result := range cn {
		results[result.key] = result
	}

	return results
}
