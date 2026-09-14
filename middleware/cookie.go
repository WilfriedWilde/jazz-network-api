package middleware

import (
	"fmt"
	"net/http"
)

func GetSessionID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_UUID")
	if err != nil {
		return "", fmt.Errorf("retrieve sessionID cookie: %w", err)
	}

	return cookie.Value, nil
}
