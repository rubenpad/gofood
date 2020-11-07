package services

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/rubbenpad/gofood/domain"
)

// Indicates if should append when exists a comma `,` or a close parenthesis `)`
func shouldAppend(ch byte) bool {
	return ch == 44 || ch == 41
}

func isOpenParenthesis(ch byte) bool {
	return ch == 40
}

func shouldChangeProcess(ch byte) bool {
	return ch == 0
}

func formatTransactionsData(
	date string,
	bytesData []byte,
	productsUids,
	buyersUids map[string]string) []domain.Transaction {

	data := bytes.Split(bytesData, []byte("\x00\x00"))
	transactions := make([]domain.Transaction, len(data)-1)

	for i := range data {
		raw := string(data[i])
		if len(raw) == 0 {
			break
		}

		var ip, device, productID string
		productsID := []domain.Uid{}

		k, process := 23, 0
		for k < len(raw) {
			ch := raw[k]

			switch {
			case process == 0 && !shouldChangeProcess(ch):
				ip += string(ch)
				k++
			case process == 1 && !shouldChangeProcess(ch):
				device += string(ch)
				k++
			case process == 2 && !shouldChangeProcess(ch):
				if isOpenParenthesis(ch) {
					k++
				} else {
					productID += string(ch)
					k++

					if shouldAppend(raw[k]) {
						p := domain.Uid{}
						if val, ok := productsUids[productID]; ok {
							p.UID = val
						} else {
							p.UID = "_:" + productID
						}

						productsID = append(productsID, p)
						productID = ""
						k++
					}
				}
			default:
				process++
				k++
			}
		}

		id := string(raw[1:13])
		buyerID := string(raw[14:22])
		owner := domain.Uid{}
		if val, ok := buyersUids[buyerID]; ok {
			owner.UID = val
		} else {
			owner.UID = "_:" + buyerID
		}

		when := domain.Timestamp{UID: "_:" + date, Date: date}
		from := domain.Ip{UID: "_:" + ip, IP: ip}

		transactions[i] = domain.Transaction{
			ID:       id,
			Device:   device,
			When:     when,
			Products: productsID,
			From:     from,
			Owner:    owner,
		}
	}

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

func formatProductsData(bytesData []byte, productsInStore []domain.Product) []domain.Product {
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

	return products
}

func formatBuyersData(data []byte, buyersInStore []domain.Buyer) []domain.Buyer {
	buyersMap := make(map[string]string, len(buyersInStore))
	for i := range buyersInStore {
		current := buyersInStore[i]
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

	return buyers
}
