package services

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/rubbenpad/gofood/domain"
)

func formatTransactionsData(
	date string,
	bytesData []byte,
	productsUids,
	buyersUids map[string]string) []byte {

	data := bytes.Split(bytesData, []byte("\x00\x00"))
	transactions := make([]domain.Transaction, len(data)-1)
	empty := []byte("")

	for i, raw := range data {
		if len(raw) == 0 {
			break
		}

		raw = bytes.Replace(raw, []byte("#"), empty, 1)
		raw = bytes.Replace(raw, []byte("("), empty, 1)
		raw = bytes.Replace(raw, []byte(")"), empty, 1)
		d := bytes.Split(raw, []byte("\x00"))

		rawProductsID := bytes.Split(d[4], []byte(","))
		productsID := make([]domain.Uid, len(rawProductsID))
		for k := range rawProductsID {
			productID := string(rawProductsID[k])

			if val, ok := productsUids[productID]; ok {
				productsID[k] = domain.Uid{UID: val}
			} else {
				productsID[k] = domain.Uid{UID: "_:" + productID}
			}
		}

		owner := domain.Uid{}
		buyerID := string(d[1])
		if val, ok := buyersUids[buyerID]; ok {
			owner.UID = val
		} else {
			owner.UID = "_:" + buyerID
		}

		ip := string(d[2])
		from := domain.Ip{UID: "_:" + ip, IP: ip}
		when := domain.Timestamp{UID: "_:" + date, Date: date}

		transactions[i] = domain.Transaction{
			ID:       string(d[0]),
			Device:   string(d[3]),
			When:     when,
			Products: productsID,
			From:     from,
			Owner:    owner,
		}
	}

	encodedTransactions, _ := json.Marshal(transactions)
	return encodedTransactions
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

func formatProductsData(bytesData []byte, productsInStore []domain.Product) []byte {
	// Creates a map to search fast what products already exists in the store
	productsMap := make(map[string]string, len(productsInStore))
	for i := range productsInStore {
		current := productsInStore[i]
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

	encodedProducts, _ := json.Marshal(products)
	return encodedProducts
}

func formatBuyersData(data []byte, buyersInStore []domain.Buyer) []byte {
	buyers := []domain.Buyer{}
	json.Unmarshal(data, &buyers)

	buyersMap := make(map[string]string, len(buyersInStore))
	for i := range buyersInStore {
		current := buyersInStore[i]
		buyersMap[current.ID] = current.UID
	}

	for k := range buyers {
		if val, ok := buyersMap[buyers[k].ID]; ok {
			buyers[k].UID = val
		} else {
			buyers[k].UID = "_:" + buyers[k].ID
		}
	}

	encodedBuyers, _ := json.Marshal(buyers)
	return encodedBuyers
}
