package routes

import (
	"encoding/json"
	"net/http"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/services"
)

func LoadDataAPI(ap *app.App) {
	dataService := services.NewDataService()

	ap.Router.Post("/data", func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")
		dataIsAlreadyLoaded, err := dataService.Load(date)

		res := response{}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if dataIsAlreadyLoaded {
			res.Message = "Data for this date is already loaded"
			res.Status = "OK"

			json.NewEncoder(w).Encode(res)
			return
		}

		if err != nil {
			res.Message = "Data no loaded"
			res.Status = "Error"

			json.NewEncoder(w).Encode(res)
			return
		}

		res.Message = "Data is being loaded"
		res.Status = "OK"
		json.NewEncoder(w).Encode(res)
		return
	})
}
