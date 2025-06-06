package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := parameters{}

	err := decoder.Decode(&reqParams)
	if err != nil {
		respondWithErr(w, http.StatusBadRequest, "Couldnt Decode parameters", err)
	}
	addedUser, err := cfg.dbQueries.CreateUser(req.Context(), reqParams.Email)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		respondWithErr(w, http.StatusInternalServerError, "Failed to create user", err)
	}
	respondWithJSON(w, http.StatusCreated, response{
		ID:        addedUser.ID,
		CreatedAt: addedUser.CreatedAt,
		UpdatedAt: addedUser.UpdatedAt,
		Email:     addedUser.Email,
	})
}
