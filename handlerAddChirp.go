package main

import (
	"encoding/json"
	"net/http"

	"github.com/hatimhas/chirtpy/internal/auth"
	"github.com/hatimhas/chirtpy/internal/database"
)

func (c *apiConfig) handlerAddChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	tokenAuth, err := auth.GetBearerToken(req.Header)
	if err != nil {

		respondWithErr(w, http.StatusUnauthorized, "Couldn't find JWT", err)

		return
	}
	userID, err := auth.ValidateJWT(tokenAuth, c.secretKey)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := parameters{}

	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldnt Decode parameters", err)
		return
	}

	if len(reqParams.Body) > 140 {
		respondWithErr(w, http.StatusBadRequest, "Chirp is too long", nil)
		return

	}

	cleanedreqBody := profaneCheck(reqParams.Body)

	chirp, err := c.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   cleanedreqBody,
		UserID: userID,
	})
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldnt create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	},
	)
}
