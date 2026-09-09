package httpapi

import (
	"jazz-network-api/model"
	"jazz-network-api/service"
	"net/http"
	"strings"
	"time"
	"uuid"
)

type SessionHandler struct {
	SessionService *service.SessionService
	UserService    *service.UserService
}

const sessionExpiryTime = time.Hour * 1

func (s *SessionHandler) Signup(w http.ResponseWriter, r *http.Request) {
	input, err := decodeJSON[model.UserInput](w, r)
	if err != nil {
		writeError(w, err)
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if !isValidEmail(input.Email) {
		writeError(w, err)
	}

	user, err := s.UserService.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
	}

	// TODO: implement session Create flow
	session, err := s.SessionService.Create(r.Context(), model.SessionCreate{
		UUID:      uuid.New().String(),
		UserID:    user.ID,
		CreatedAt: time.Now().Format(time.DateTime),
		ExpiresAt: time.Now().Add(sessionExpiryTime).Format(time.DateTime),
	})
	if err != nil {
		writeError(w, err)
	}

	// TODO: set cookie
}
