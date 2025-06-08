package main

import (
	"net/http"
	"time"

	"github.com/hatimhas/chirtpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Extract token from Authorization header
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Missing or malformed Authorization header", err)
		return
	}

	// Fetch user via refresh token
	refreshTokenRow, err := cfg.dbQueries.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	expiresIn := time.Hour
	token, err := auth.MakeJWT(refreshTokenRow.UserID, cfg.secretKey, expiresIn)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Failed to generate JWT", err)
		return
	}

	// Respond with the new token
	respondWithJSON(w, http.StatusOK, response{
		Token: token,
	})
}
