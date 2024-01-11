package middlewares

import (
	"context"
	"movies-auth/users/internal/domain"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type SessionService interface {
	ValidateSession(key uuid.UUID) (domain.Session, error)
}

type sessionKey string

var SessionKey sessionKey = "sessionKey"

func Auth(ss SessionService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "login") || strings.Contains(r.URL.Path, "register") || strings.Contains(r.URL.Path, "sessions") {
				next.ServeHTTP(w, r)
			}

			sessionCookie, err := r.Cookie("session")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			sessionKey, err := uuid.Parse(sessionCookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			session, err := ss.ValidateSession(sessionKey)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), SessionKey, session.Key)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
