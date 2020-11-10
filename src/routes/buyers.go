package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/services"
)

func BuyersAPI(ap *app.App) {
	buyersService := services.NewBuyersService()

	ap.Router.Get("/buyers", func(w http.ResponseWriter, r *http.Request) {
		data, err := buyersService.FindAllBuyers()
		res := response{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err != nil {
			res.Message = "Failed to fetch buyer's data"
			res.Status = "Error"
			json.NewEncoder(w).Encode(res)
			return
		}

		buyers := decodeBuyers{}
		json.Unmarshal(data, &buyers)

		res.Data = buyers
		res.Message = "Success"
		res.Status = "OK"
		json.NewEncoder(w).Encode(res)
		return
	})

	ap.Router.Get("/buyers/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		data, err := buyersService.FindTransactions(id)

		res := response{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err != nil {
			res.Message = "Failed to fetch buyer's data"
			res.Status = "Error"
			json.NewEncoder(w).Encode(res)
			return
		}

		datadecoded := decodeddata{}
		json.Unmarshal(data, &datadecoded)

		if len(datadecoded.Buyer) == 0 {
			w.WriteHeader(http.StatusNotFound)
			res.Message = "Not Found"
			res.Status = "OK"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Data = datadecoded
		res.Message = "Success"
		res.Status = "OK"
		json.NewEncoder(w).Encode(res)
		return
	})
}
