package routes

import (
	"encoding/json"
	"net/http"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/services"
)

func LoadDataAPI(ap *app.App) {
	loadDataService := services.NewloadDataService()

	ap.Router.Get("/data", func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")
		err := loadDataService.GetData(date)
		res := response{}
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			res.Data = make([]interface{}, 1)
			res.Message = "Data no loaded"
			res.Status = "Error"

			json.NewEncoder(w).Encode(res)
		}

		res.Data = nil
		res.Message = "Data loaded"
		res.Status = "OK"

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	})
}
