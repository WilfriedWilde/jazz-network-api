package middleware

import (
	"context"
	"jazz-network-api/model"
	repo "jazz-network-api/repository"
	"net/http"
)

type contextKey string

var userContextKey contextKey = "user"

func UserFromContext(ctx context.Context) (model.User, bool) {
	user, ok := ctx.Value(userContextKey).(model.User)
	return user, ok
}

func AuthMiddleware(sessionRepo *repo.SessionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionUUID, err := GetSessionID(r)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, model.ErrBody{
					Code:    "unauthorized",
					Message: "session not found",
				})
				return
			}

			user, err := sessionRepo.GetUserBySession(r.Context(), sessionRepo.DB, sessionUUID)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, model.ErrBody{
					Code:    "unauthorized",
					Message: "user not found",
				})
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
