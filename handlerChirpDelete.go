package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hatimhas/chirtpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithErr(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}
	if dbChirp.UserID != userID {
		respondWithErr(w, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
