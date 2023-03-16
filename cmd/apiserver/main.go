package main

import (
	"log"
	"net/http"

	"github.com/microsoft/commercial-marketplace-offer-deploy/cmd/apiserver/routes"
)

func main() {
	log.Printf("Server started")

	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
