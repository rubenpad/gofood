package services

import (
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/rubbenpad/gofood/domain"
)

// Endpoint "/transactions" send no standard data
// i.e: "#00005f80fa12'2a2dc5b'246.124.213.49'ios'(7dd44f1d,e4356fea)"
// representing data from transactions and this function is a helper to format it.
func formatNoStandardData(content io.ReadCloser) []domain.Transaction {
	data, err := ioutil.ReadAll(content)
	if err != nil {
		log.Println("Couldn't format data")
	}

	raw := strings.Split(string(data), "\x00\x00")
	transactions := make([]domain.Transaction, len(raw)-1)

	for i := 0; i < len(raw)-1; i++ {
		str := raw[i]
		str = strings.Replace(str, "#", "", -1)
		str = strings.Replace(str, "(", "", -1)
		str = strings.Replace(str, ")", "", -1)

		rawstr := strings.Split(str, "\x00")
		productIdsRaw := strings.Split(rawstr[4], ",")
		productsid := make([]domain.Uid, len(productIdsRaw))

		for j, pid := range productIdsRaw {
			productsid[j] = domain.Uid{UID: "_:" + pid}
		}

		transactions[i] = domain.Transaction{
			ID:       rawstr[0],
			Device:   rawstr[3],
			Products: productsid,
			From:     domain.Ip{IP: rawstr[2], UID: "_:" + rawstr[2]},
			Owner:    domain.Uid{UID: "_:" + rawstr[1]},
		}
	}

	return transactions
}

// /products data is formatted like CSV but with ' as separator
// This function returns received data as an slice of
// { "id": product_id, "name": product_name, "price": product_price }
func formatCSVData(content io.ReadCloser) []domain.Product {
	data, err := ioutil.ReadAll(content)
	if err != nil {
		log.Println("Couldn't format data")
	}

	raw := strings.Split(string(data), "\n")
	products := make([]domain.Product, len(raw)-1)
	leftRegex := regexp.MustCompile(`(?P<left>[a-z0-9])(?:')(?P<right>[A-Z])`)
	rightRegex := regexp.MustCompile(`(?P<left>[a-z])(?:')(?P<right>[0-9])`)

	for i := 0; i < len(raw)-1; i++ {
		// Work to format data. Here delete double quote and replace the leftmost
		// and rightmost single quote by a comma then split current item to crate
		// product struct and append it to products slice.
		item := raw[i]
		item = strings.Replace(item, "\"", "", -1)
		item = leftRegex.ReplaceAllString(item, "$left,$right")
		item = rightRegex.ReplaceAllString(item, "$left,$right")
		rawItem := strings.Split(item, ",")
		products[i] = domain.Product{
			UID:   "_:" + rawItem[0],
			Name:  rawItem[1],
			Price: rawItem[2],
		}
	}

	return products
}
