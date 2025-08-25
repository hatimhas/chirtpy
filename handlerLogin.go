package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hatimhas/chirtpy/internal/auth"
	"github.com/hatimhas/chirtpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		RefreshToken string `json:"refresh_token"`
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

	expiresIn := time.Duration(3600) * time.Second
	accessToken, err := auth.MakeJWT(userDB.ID, cfg.secretKey, expiresIn)
	if err != nil {
		fmt.Printf("Failed to create JWT token: %v\n", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Println(err)
		return
	}

	AddRefTokenParams, err := cfg.dbQueries.AddRefToken(req.Context(), database.AddRefTokenParams{
		Token:  refreshToken,
		UserID: userDB.ID,
	})
	if err != nil {
		fmt.Printf("Failed to add refresh token to db: %v\n", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        userDB.ID,
			CreatedAt: userDB.CreatedAt,
			UpdatedAt: userDB.UpdatedAt,
			Email:     userDB.Email,
			Token:     accessToken,
			ChirpyRed: userDB.IsChirpyRed,
		},
		RefreshToken: AddRefTokenParams.Token,
	})
}
