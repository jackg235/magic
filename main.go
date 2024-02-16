package main

import (
	"log"
	"magic-ledger/api"
	"net/http"
)

func main() {
	api.InitializeInternalAccounts()
	router := api.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))

}
