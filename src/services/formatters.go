package services

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/rubbenpad/gofood/domain"
)

func formatTransactionsData(date string, data []byte, productsUids, buyersUids map[string]string) []domain.Transaction {
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
			if val, ok := productsUids[pid]; ok {
				productsid[j] = domain.Uid{UID: val}
			} else {
				productsid[j] = domain.Uid{UID: "_:" + pid}
			}
		}

		owner := domain.Uid{}
		if val, ok := buyersUids[rawstr[1]]; ok {
			owner.UID = val
		} else {
			owner.UID = "_:" + rawstr[1]
		}

		when := domain.Timestamp{UID: "_:" + date, Date: date}
		transactions[i] = domain.Transaction{
			ID:       rawstr[0],
			Device:   rawstr[3],
			When:     when,
			Products: productsid,
			From:     domain.Ip{IP: rawstr[2], UID: "_:" + rawstr[2]},
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
