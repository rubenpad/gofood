package routes

import (
	"net/http"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/services"
)

func LoadDataAPI(ap *app.App) {

	loadDataService := services.NewLoadDataService()

	ap.Router.Get("/load", func(w http.ResponseWriter, r *http.Request) {
		loadDataService.GetData()
		w.Write([]byte("loading..."))
	})
}
