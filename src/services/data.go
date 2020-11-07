package services

import (
	"encoding/json"
	"io/ioutil"
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
	if dateExists := store.GetDate(date); dateExists {
		return dateExists, nil
	}

	// Build requests to remote data and fetch concurrently
	requests := ld.buildRequests(date)
	results := ld.fetchConcurrently(requests)

	// Format, encode and save products and buyers data
	all := store.FindAll()
	products := formatProductsData(results["products"].response.data, all.Products)
	buyers := formatBuyersData(results["buyers"].response.data, all.Buyers)

	encodedProducts, _ := json.Marshal(products)
	encodedBuyers, _ := json.Marshal(buyers)
	assignedProducts, _ := store.Save(encodedProducts)
	assignedBuyers, _ := store.Save(encodedBuyers)

	// Format, encode and save transactions data
	transactions := formatTransactionsData(
		date,
		results["transactions"].response.data,
		assignedProducts.Uids,
		assignedBuyers.Uids,
	)
	encodedTransactions, _ := json.Marshal(transactions)
	if _, err := store.Save(encodedTransactions); err != nil {
		return false, err
	}

	return false, nil
}

func (ld *loadDataService) buildRequests(date string) map[string]func() (*remoteResponse, error) {
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
			return ld.makeRequest(baseurl + endpoint + date)
		}
	}

	return requests
}

type remoteResponse struct {
	data []byte
}

func (ld *loadDataService) makeRequest(url string) (*remoteResponse, error) {
	res, err := ld.httpclient.Get(url)
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

func (ld *loadDataService) fetchConcurrently(requests map[string]func() (*remoteResponse, error)) map[string]*requestResult {
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
