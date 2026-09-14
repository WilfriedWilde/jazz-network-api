package httpapi

import (
	"errors"
	"jazz-network-api/model"
	"jazz-network-api/service"
	"jazz-network-api/validation"
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
	if !validation.IsValidEmail(input.Email) {
		writeError(w, &model.RequestError{
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
	if !validation.IsValidEmail(input.Email) {
		writeError(w, &model.RequestError{
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

func (s *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// get sessionUUID
	cookie, err := r.Cookie("session_UUID")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	// delete session
	if err := s.SessionService.Delete(r.Context(), cookie.Value); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeJSON(w, http.StatusUnauthorized, model.ErrBody{
				Code:    "unauthorized",
				Message: "session not found",
			})
		} else {
			writeError(w, err)
		}

		return
	}

	// clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_UUID",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	// write successful response
	w.WriteHeader(http.StatusOK)
}
