package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hatimhas/chirtpy/internal/auth"
)

// TODO Query lookup using email body in req, compare the pw in db with pw in req. If lookup of email/matching pw fail return 401 unauthorized msg "Incorrect email or password". If match return 200 ok with JSON response (no PW).
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:        userDB.ID,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
		Email:     userDB.Email,
	})
}
