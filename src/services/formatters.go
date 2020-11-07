package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/rubbenpad/gofood/domain"
)

func formatTransactionsData(
	date string,
	bytesData []byte,
	productsUids,
	buyersUids map[string]string) []domain.Transaction {

	ts := time.Now()
	splittedData := bytes.Split(bytesData, []byte("\x00\x00"))
	transactions := make([]domain.Transaction, len(splittedData)-1)

	for i, raw := range splittedData {
		if len(raw) == 0 {
			break
		}

		raw = bytes.Replace(raw, []byte("#"), []byte(""), -1)
		raw = bytes.Replace(raw, []byte("("), []byte(""), -1)
		raw = bytes.Replace(raw, []byte(")"), []byte(""), -1)

		data := bytes.Split(raw, []byte("\x00"))
		rawProductsID := bytes.Split(data[4], []byte(","))
		productsid := make([]domain.Uid, len(rawProductsID))

		for j, v := range rawProductsID {
			pid := string(v)
			if val, ok := productsUids[pid]; ok {
				productsid[j] = domain.Uid{UID: val}
			} else {
				productsid[j] = domain.Uid{UID: "_:" + pid}
			}
		}

		owner := domain.Uid{}
		buyerID := string(data[1])
		if val, ok := buyersUids[buyerID]; ok {
			owner.UID = val
		} else {
			owner.UID = "_:" + buyerID
		}

		when := domain.Timestamp{UID: "_:" + date, Date: date}
		from := domain.Ip{IP: string(data[2]), UID: "_:" + string(data[2])}

		transactions[i] = domain.Transaction{
			ID:       string(data[0]),
			Device:   string(data[3]),
			When:     when,
			Products: productsid,
			From:     from,
			Owner:    owner,
		}
	}

	te := time.Now()
	fmt.Println("Transactions time: ", te.Sub(ts))
	return transactions
}

// Helper functions to format products data
func isSingleQuote(ch byte) bool {
	return ch == 39
}

func isDoubleQuote(ch byte) bool {
	return ch == 34
}

func isDigit(ch byte) bool {
	return ch >= 48 && ch <= 57
}

type productsInStore struct {
	Products []domain.Product `json:"products"`
}

func formatProductsData(bytesData, savedProducts []byte) []domain.Product {
	ts := time.Now()
	sip := productsInStore{}
	sipError := json.Unmarshal(savedProducts, &sip)
	if sipError != nil {
		log.Panic("Error trying to decode data")
	}

	// Creates a map to search fast what products already exists in the store
	productsMap := make(map[string]string, len(sip.Products))
	for i := range sip.Products {
		current := sip.Products[i]
		productsMap[current.ID] = current.UID
	}

	data := bytes.Split(bytesData, []byte("\n"))
	products := make([]domain.Product, len(data)-1)

	for i := range data {
		raw := string(data[i])
		if len(raw) == 0 {
			break
		}

		var k int
		var id string
		if isSingleQuote(raw[7]) {
			id = raw[:7]
			k = 8
		} else {
			id = raw[:8]
			k = 9
		}

		start, end := k, 0
		quoted := false
		for k < len(raw) {
			if isDoubleQuote(raw[start]) {
				quoted = true
				start++
				k++
			} else if isSingleQuote(raw[k]) && isDigit(raw[k+1]) {
				if quoted {
					end = k - 1
					k++
					break
				} else {
					end = k
					k++
					break
				}
			}
			k++
		}

		// Assign uid and price
		var uid string
		rawPrice, _ := strconv.Atoi(raw[k:len(raw)])
		if val, ok := productsMap[id]; ok {
			uid = val
		} else {
			uid = "_:" + id
		}

		products[i] = domain.Product{
			UID:   uid,
			ID:    id,
			Name:  raw[start:end],
			Price: rawPrice,
		}
	}

	te := time.Now()
	fmt.Println("Products time: ", te.Sub(ts))
	return products
}

type sb struct {
	Buyers []domain.Buyer `json:"buyers"`
}

func formatBuyersData(data, savedBuyers []byte) []domain.Buyer {
	ts := time.Now()
	sby := sb{}
	sbE := json.Unmarshal(savedBuyers, &sby)
	if sbE != nil {
		log.Panic("Error trying to decode data")
	}

	buyersMap := make(map[string]string, len(sby.Buyers))
	for i := range sby.Buyers {
		current := sby.Buyers[i]
		buyersMap[current.ID] = current.UID
	}

	buyers := []domain.Buyer{}
	json.Unmarshal(data, &buyers)

	for k := range buyers {
		if val, ok := buyersMap[buyers[k].ID]; ok {
			buyers[k].UID = val
		} else {
			buyers[k].UID = "_:" + buyers[k].ID
		}
	}

	te := time.Now()
	fmt.Println("Transactions time: ", te.Sub(ts))
	return buyers
}
