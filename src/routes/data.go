package routes

import (
	"encoding/json"
	"net/http"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/services"
)

func LoadDataAPI(ap *app.App) {
	loadDataService := services.NewloadDataService()

	ap.Router.Post("/data", func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")
		dataIsAlreadyLoaded, err := loadDataService.GetData(date)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res := response{}

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

		res.Message = "Data loaded"
		res.Status = "OK"
		json.NewEncoder(w).Encode(res)
		return
	})
}
