package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/services"
)

func LoadDataAPI(ap *app.App) {

	loadDataService := services.NewloadDataService()

	ap.Router.Get("/load", func(w http.ResponseWriter, r *http.Request) {
		query := loadDataService.GetData()
		encoded, _ := json.Marshal(query)
		fmt.Fprintf(w, "%s", encoded)
	})
}
