package middlewares

import (
	"context"

	"go-rest-template/internal/db"
	"go-rest-template/internal/models"
	"go-rest-template/pkg"
	"go-rest-template/pkg/api"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user"

// Middleware that checks if the user is authenticated and adds the user to the contextKey "user" if authenticated
func IsAuthenticated(q *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				api.Unauthorized(w, "Unauthorized")
				return
			}

			token := cookie.Value
			email, err := pkg.ValidateToken(token)
			if err != nil {
				api.Unauthorized(w, "Unauthorized")
				return
			}

			user, err := q.GetUserByEmail(r.Context(), email)
			if err != nil || user.ID == 0 {
				api.Unauthorized(w, "Unauthorized")
				return
			}
			ctx := context.WithValue(r.Context(), userContextKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// returns the user from the context
func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value(userContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
