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
		err := loadDataService.GetData()
		if err != nil {
			fmt.Fprint(w, "Service couldn't load data")
		}
		fmt.Fprint(w, "Data loaded successfully")
	})
}
