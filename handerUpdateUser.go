package main

import (
	"encoding/json"
	"net/http"

	"github.com/hatimhas/chirtpy/internal/auth"
	"github.com/hatimhas/chirtpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}
	accesssToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Missing or malformed Authorization header", err)
		return
	}
	userID, err := auth.ValidateJWT(accesssToken, cfg.secretKey)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := parameters{}

	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Couldnt Decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user, err := cfg.dbQueries.UpdateUser(req.Context(), database.UpdateUserParams{
		Email:          reqParams.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
