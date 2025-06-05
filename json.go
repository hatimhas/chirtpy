package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithErr(w http.ResponseWriter, status int, msg string, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	if err != nil {
		log.Printf("Error: %s", err)
	}
	if status > 499 {
		log.Printf("Server error: %s", msg)
	}
	respondWithJSON(w, status, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	w.Write(append(response, '\n'))
}
