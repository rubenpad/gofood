package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rubbenpad/gofood/app"
	"github.com/rubbenpad/gofood/routes"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	ap := app.New()
	routes.LoadDataAPI(ap)

	log.Fatal(http.ListenAndServe(":3000", ap))
}
