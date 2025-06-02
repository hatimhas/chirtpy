package main

import (
	"log"
	"net/http"
)

func main() {
	const assetsDir = "./assets"
	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	// Use Handle method of serveMux
	serveMux.Handle("/", http.FileServer(http.Dir("./static")))
	serveMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	log.Fatal(server.ListenAndServe())
}
