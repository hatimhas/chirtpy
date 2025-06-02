package main

import (
	"log"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	// Use Handle method of serveMux
	serveMux.Handle("/", http.FileServer(http.Dir("./static")))

	log.Fatal(server.ListenAndServe())
}
