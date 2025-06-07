package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hatimhas/chirtpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
	}
	decoder := json.NewDecoder(req.Body)
	reqParams := parameters{}

	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Couldnt Decode parameters", err)
		return
	}

	userDB, err := cfg.dbQueries.GetUserByEmail(req.Context(), reqParams.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithErr(w, http.StatusUnauthorized, "Incorrect email or password ", nil)
			return
		}
		respondWithErr(w, http.StatusInternalServerError, "Failed to retrieve user", err)
	}

	pwMatch := auth.CheckPasswordHash(userDB.HashedPassword, reqParams.Password)
	if pwMatch != nil {
		respondWithErr(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	expTime := reqParams.ExpiresInSeconds
	if expTime > 3600 || expTime == 0 {
		expTime = 3600
	}
	expiresIn := time.Duration(expTime) * time.Second
	accessToken, err := auth.MakeJWT(userDB.ID, cfg.secretKey, expiresIn)
	if err != nil {
		fmt.Printf("Failed to create JWT token: %v\n", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        userDB.ID,
			CreatedAt: userDB.CreatedAt,
			UpdatedAt: userDB.UpdatedAt,
			Email:     userDB.Email,
			Token:     accessToken,
		},
	})
}
