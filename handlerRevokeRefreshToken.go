package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/hatimhas/chirtpy/internal/auth"
	"github.com/hatimhas/chirtpy/internal/database"
)

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, req *http.Request) {
	// Get token from Authorization header
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		http.Error(w, "Invalid or missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Revoke token in DB: set revoked_at and updated_at to now
	now := time.Now().UTC()
	err = cfg.dbQueries.RevokeRefreshToken(req.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: now, Valid: true},
		UpdatedAt: now,
		Token:     refreshToken,
	})
	if err != nil {
		http.Error(w, "Failed to revoke token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}
