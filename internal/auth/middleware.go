package auth

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const UserContextKey = contextKey("user")

type User struct {
	ID       string
	Username string
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple authentication - in production, use JWT or sessions
		username := strings.TrimSpace(r.URL.Query().Get("username"))
		if username == "" {
			http.Error(w, "Username is required", http.StatusUnauthorized)
			return
		}

		user := &User{
			ID:       generateUserID(),
			Username: username,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateUserID() string {
	return "user_" + string(time.Now().UnixNano())
}
