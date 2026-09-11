package httpapi

import (
	"jazz-network-api/model"
	"jazz-network-api/service"
	"net/http"
	"strings"
	"time"
)

type SessionHandler struct {
	SessionService *service.SessionService
	UserService    *service.UserService
}

func (s *SessionHandler) Signup(w http.ResponseWriter, r *http.Request) {
	input, err := decodeJSON[model.UserInput](w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if !isValidEmail(input.Email) {
		writeError(w, &requestError{
			Status:  http.StatusBadRequest,
			Code:    "validation_error",
			Message: "invalid email format",
		})
		return
	}

	user, err := s.UserService.Create(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	session, err := s.SessionService.Create(r.Context(), user.ID)
	if err != nil {
		writeError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_UUID",
		Value:    session.UUID,
		Expires:  time.Now().Add(service.SessionExpiryTime),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	writeJSON(w, http.StatusOK, user)
}

func (s *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {
	input, err := decodeJSON[model.UserInput](w, r)
	if err != nil {
		writeError(w, err)
		return
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if !isValidEmail(input.Email) {
		writeError(w, &requestError{
			Status:  http.StatusBadRequest,
			Code:    "validation_error",
			Message: "invalid email format",
		})
		return
	}

	user, err := s.UserService.Authenticate(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	session, err := s.SessionService.Create(r.Context(), user.ID)
	if err != nil {
		writeError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_UUID",
		Value:    session.UUID,
		Expires:  time.Now().Add(service.SessionExpiryTime),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	writeJSON(w, http.StatusOK, user)
}
