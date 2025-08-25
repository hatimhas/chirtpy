package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hatimhas/chirtpy/internal/auth"
)

func (c *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiKey(req.Header)
	if err != nil {
		respondWithErr(w, http.StatusUnauthorized, "Missing or malformed Authorization header", err)
		return
	}
	if apiKey != c.polkaKey {
		respondWithErr(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := parameters{}

	err = decoder.Decode(&reqParams)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldnt Decode parameters", err)
		return
	}

	if reqParams.Event != "user.upgraded" {
		respondWithErr(w, http.StatusNoContent, "Invalid Webhook Event !user.upgraded", nil)
		return

	}

	userID, err := uuid.Parse(reqParams.Data.UserID)
	if err != nil {
		fmt.Println("Invalid UUID in params:", err)
		return
	}
	_, err = c.dbQueries.UpdateUserChirpyRed(req.Context(), userID)
	if err != nil {
		respondWithErr(w, http.StatusInternalServerError, "Couldnt Find UserID in DB", err)

		return
	}

	respondWithJSON(w, http.StatusNoContent, parameters{})
}
