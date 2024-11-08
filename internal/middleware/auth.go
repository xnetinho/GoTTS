package middleware

import (
	"net/http"
	"tts-api/internal/handlers"
)

func AuthMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			if authToken != "Bearer "+token {
				handlers.WriteJSONError(w, http.StatusUnauthorized, "Não autorizado")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
