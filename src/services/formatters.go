package services

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/rubbenpad/gofood/domain"
)

func formatTransactionsData(
	date string,
	bytesData []byte,
	productsUids,
	buyersUids map[string]string) []domain.Transaction {

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
	return transactions
}

type sp struct {
	Products []domain.Product `json:"products"`
}

func formatProductsData(data, savedProducts []byte) []domain.Product {
	sP := sp{}
	spE := json.Unmarshal(savedProducts, &sP)
	if spE != nil {
		log.Panic("Error trying to decode data")
	}

	productsMap := make(map[string]string, len(sP.Products))
	for i := range sP.Products {
		current := sP.Products[i]
		productsMap[current.ID] = current.UID
	}

	raw := strings.Split(string(data), "\n")
	products := make([]domain.Product, len(raw)-1)
	regex := regexp.MustCompile(`(?P<left>[\w\W])(?:')(?P<right>[0-9])`)

	for i := 0; i < len(raw)-1; i++ {
		item := raw[i]
		item = strings.Replace(item, "\"", "", -1)
		item = strings.Replace(item, "'", ",", 1)
		item = regex.ReplaceAllString(item, "$left,$right")
		rawItem := strings.Split(item, ",")
		price, _ := strconv.Atoi(rawItem[2])

		uid := ""
		if val, ok := productsMap[rawItem[0]]; ok {
			uid = val
		} else {
			uid = "_:" + rawItem[0]
		}

		products[i] = domain.Product{
			UID:   uid,
			ID:    rawItem[0],
			Name:  rawItem[1],
			Price: price,
		}
	}

	return products
}

type sb struct {
	Buyers []domain.Buyer `json:"buyers"`
}

func formatBuyersData(data, savedBuyers []byte) []domain.Buyer {
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

	return buyers
}
