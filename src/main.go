package main

import (
	"log"
	"net/http"

	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/routes"
)

func main() {

	ap := app.New()
	routes.LoadDataAPI(ap)

	log.Fatal(http.ListenAndServe(":3000", ap))
}
