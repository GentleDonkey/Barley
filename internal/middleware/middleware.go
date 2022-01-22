package middleware

import (
	"net/http"
	"notifications/internal/jwt"
)

func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if jwt.TokenParse(w, r) {
			next.ServeHTTP(w, r)
		}
		return
	})
}
