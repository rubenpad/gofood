package routes

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/services"
)

func BuyersAPI(ap *app.App) {
	buyersService := services.NewBuyersService()

	ap.Router.Get("/buyers/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		buyersService.FindTransactions(id)
		/*if err != nil {
			jsonapi.MarshalPayload(w, response{
				data:    nil,
				message: "Error trying to fetch buyer transactions",
				status:  "Error",
			})
		}
		jsonapi.MarshalPayload(w, response{
			data:    data,
			message: "Success",
			status:  "OK",
		})*/
	})
}
