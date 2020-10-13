package routes

import (
	"fmt"
	"net/http"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/services"
)

func LoadDataAPI(ap *app.App) {

	loadDataService := services.NewloadDataService()

	ap.Router.Get("/load", func(w http.ResponseWriter, r *http.Request) {
		transactions := loadDataService.GetData()
		fmt.Fprint(w, transactions)
	})
}
