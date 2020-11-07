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

	// Build a hash map of request to then fetch data concurrently
	requests := ld.buildRequests(date)
	results := ld.fetchConcurrently(requests)

	// Get data already saved `products` and `buyers` to compare them with new data
	// and avoid have duplicates.
	all := store.FindAll()
	products := formatProductsData(results["products"].response.data, all.Products)
	buyers := formatBuyersData(results["buyers"].response.data, all.Buyers)

	// First encode and save products and buyers then transactions.
	encodedProducts, _ := json.Marshal(products)
	assignedProducts, _ := store.Save(encodedProducts)

	encodedBuyers, _ := json.Marshal(buyers)
	assignedBuyers, _ := store.Save(encodedBuyers)

	transactions := formatTransactionsData(
		date,
		results["transactions"].response.data,
		assignedProducts.Uids,
		assignedBuyers.Uids,
	)
	encodedTransactions, _ := json.Marshal(transactions)
	_, err := store.Save(encodedTransactions)

	if err != nil {
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
