package auth

import (
	"fmt"
	"net/http"
)

func getSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", fmt.Errorf("retrieve sessionID cookie: %w", err)
	}

	return cookie.Value, nil
}
