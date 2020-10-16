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

	ap.Router.Get("/buyers/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		data, err := buyersService.FindTransactions(id)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := response{}

		if err != nil {
			res.Message = "Failed to fetch buyer's data"
			res.Status = "Error"
			json.NewEncoder(w).Encode(res)
			return
		}

		datadecoded := decodeddata{}
		json.Unmarshal(data, &datadecoded)

		res.Data = datadecoded
		res.Message = "Success"
		res.Status = "OK"
		json.NewEncoder(w).Encode(res)
		return
	})
}
