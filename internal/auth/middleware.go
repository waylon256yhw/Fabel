package auth

import (
	"context"
	"net/http"

	"fabel/internal/dbq"

	"github.com/alexedwards/scs/v2"
)

type contextKey string

const userKey contextKey = "user"

// UserFromContext retrieves the authenticated user from context.
// Returns nil if no user is set (unauthenticated).
func UserFromContext(ctx context.Context) *dbq.User {
	u, _ := ctx.Value(userKey).(*dbq.User)
	return u
}

// RequireAuth is middleware that checks the session for a user_id,
// loads the user from DB, and stores it in context. Returns 401 if not authenticated.
func RequireAuth(sm *scs.SessionManager, q *dbq.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := sm.GetString(r.Context(), "user_id")
			if userID == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := q.GetUserByID(r.Context(), userID)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin is middleware that checks the context user has admin role.
// Must be used after RequireAuth.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user == nil || user.Role != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
