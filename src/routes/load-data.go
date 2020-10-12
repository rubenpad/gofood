package routes

import (
	"net/http"

	"github.com/rubbenpad/gofood/app"
)

func LoadDataAPI(ap *app.App) {
	ap.Router.Get("/load", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("loading..."))
	})
}
