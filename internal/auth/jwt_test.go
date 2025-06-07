package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "mysecretkey"
	userID := uuid.New()

	t.Run("Valid token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Minute)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		parsedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("failed to validate JWT: %v", err)
		}

		if parsedID != userID {
			t.Errorf("expected userID %v, got %v", userID, parsedID)
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Minute) // already expired
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Errorf("expected error for expired token, got none")
		}
	})

	t.Run("Wrong secret", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Minute)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		_, err = ValidateJWT(token, "wrongsecret")
		if err == nil {
			t.Errorf("expected error for wrong secret, got none")
		}
	})

	t.Run("Malformed token", func(t *testing.T) {
		_, err := ValidateJWT("not.a.jwt", secret)
		if err == nil {
			t.Errorf("expected error for malformed token, got none")
		}
	})
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
