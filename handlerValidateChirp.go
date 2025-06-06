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
		Cleaned_body string `json:"cleaned_body"`
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

	cleanedreqBody := profaneCheck(reqParams.Body)

	respondWithJSON(w, http.StatusOK, validResponse{
		Cleaned_body: cleanedreqBody,
	},
	)
}
