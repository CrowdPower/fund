package main

import (
	"log"
	"net/http"

	"github.com/crowdpower/fund/server"
)

func main() {
	r := server.Router()
	port := ":8080"
	cert := "server.crt"
	key := "server.key"
	log.Printf("Starting server on %v", port)
	log.Fatal(http.ListenAndServeTLS(port, cert, key, r))
}
