package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hatimhas/chirtpy/internal/auth"
	"github.com/hatimhas/chirtpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	hashedPassword, err := auth.HashPassword(reqParams.Password)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	addedUser, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          reqParams.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("failed to create user: %v", err)
		respondWithErr(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        addedUser.ID,
			CreatedAt: addedUser.CreatedAt,
			UpdatedAt: addedUser.UpdatedAt,
			Email:     addedUser.Email,
		},
	})
}
