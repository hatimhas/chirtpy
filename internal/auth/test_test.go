package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    http.Header
		wantToken string
		expectErr bool
	}{
		{
			name: "valid token",
			header: http.Header{
				"Authorization": []string{"Bearer abc123"},
			},
			wantToken: "abc123",
			expectErr: false,
		},
		{
			name:      "no header",
			header:    http.Header{},
			expectErr: true,
		},
		{
			name: "wrong prefix",
			header: http.Header{
				"Authorization": []string{"Token abc123"},
			},
			expectErr: true,
		},
		{
			name: "empty token",
			header: http.Header{
				"Authorization": []string{"Bearer   "},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("got unexpected error: %v", err)
				}
				if token != tt.wantToken {
					t.Errorf("expected token %q, got %q", tt.wantToken, token)
				}
			}
		})
	}
}
