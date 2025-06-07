package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	// check prefix
	prefix := "Bearer "

	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("authorization header does not start with %q", prefix)
	}

	// use len(prefix): and not TrimPrefix (more efficient), since prefix fixed lenght and start beggining of authHeader
	token := strings.TrimSpace(authHeader[len(prefix):])
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}
