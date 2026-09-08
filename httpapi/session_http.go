package httpapi

import (
  "net/http"
	"jazz-network-api/service"
)

type SessionHandler struct {
	Service *service.SessionService
}

func (s *SessionHandler) Signup (w http.ResponseWriter, r *http.Request) {
  credentials, err := parseCredentials(r)
  if err != nil {
    writeError(w, err)
  }

  user, err := service.CreateUser(r.Context(), credentials)
  if err != nil {
    writeError(w, err)
  }

  uuid := 
  session, err := service.CreateSession(r.Context(), models.SessionCreate)
}
