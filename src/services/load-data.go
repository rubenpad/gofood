package services

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const baseURL = "https://kqxty15mpg.execute-api.us-east-1.amazonaws.com"

type loadDataService struct {
	url        string
	httpclient *http.Client
}

type uid struct {
	UID string `json:"uid,omitempty"`
}

type ip struct {
	IP  string `json:"ip,omitempty"`
	UID string `json:"uid,omitempty"`
}

type transaction struct {
	ID       string `json:"id,omitempty"`
	Device   string `json:"device,omitempty"`
	Products []uid  `json:"products_id,omitempty"`
	From     ip     `json:"from,omitempty"`
	Owner    uid    `json:"owner,omitempty"`
}

func NewloadDataService() *loadDataService {
	return &loadDataService{url: baseURL, httpclient: &http.Client{}}
}

func (ld *loadDataService) GetData() []transaction {
	res, err := ld.makeRequest("/transactions?date=" + "1602530864")
	if err != nil {
		log.Fatal(err)
	}

	transactions := formatNoStandardData(res.Body)
	return transactions
}

// makeRequest is useful to perform http request to external
// endpoints to bring data
func (ld *loadDataService) makeRequest(path string) (*http.Response, error) {
	res, err := ld.httpclient.Get(ld.url + path)
	if err != nil {
		log.Fatal("Failed fetching data")
		return nil, err
	}

	return res, nil
}

// Endpoint "/transactions" send no standard data
// i.e: "#00005f80fa12'2a2dc5b'246.124.213.49'ios'(7dd44f1d,e4356fea)"
// representing data from transactions and this function is a helper to format it.
func formatNoStandardData(content io.ReadCloser) []transaction {
	data, err := ioutil.ReadAll(content)
	if err != nil {
		log.Println("Couldn't format data")
	}

	raw := strings.Split(string(data), "\x00\x00")
	transactions := make([]transaction, len(raw)-1)

	for i := 0; i < len(raw)-1; i++ {
		str := raw[i]
		str = strings.Replace(str, "#", "", -1)
		str = strings.Replace(str, "(", "", -1)
		str = strings.Replace(str, ")", "", -1)

		rawstr := strings.Split(str, "\x00")
		productIdsRaw := strings.Split(rawstr[4], ",")
		productsid := make([]uid, len(productIdsRaw))

		for j, pid := range productIdsRaw {
			productsid[j] = uid{UID: "_:" + pid}
		}

		transactions[i] = transaction{
			ID:       rawstr[0],
			Device:   rawstr[3],
			Products: productsid,
			From:     ip{IP: rawstr[2], UID: "_:" + rawstr[2]},
			Owner:    uid{UID: "_:" + rawstr[1]},
		}
	}

	return transactions
}
