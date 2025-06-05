package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type reqParameters struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(req.Body)
	reqParams := reqParameters{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldnt Decode parameters", err)
	}

	if len(reqParams.Body) > 140 {
		respondWithErr(w, http.StatusBadRequest, "Chirp is too long", nil)
		return

	}

	respondWithJSON(w, http.StatusOK, validResponse{
		Valid: true,
	},
	)
}
