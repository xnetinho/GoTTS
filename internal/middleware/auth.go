package middleware

import (
	"net/http"
	"tts-api/internal/handlers"
)

func AuthMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Rotas públicas que não exigem autenticação
			publicPaths := []string{
				"/healthcheck",
				"/api/",
			}

			for _, path := range publicPaths {
				if r.URL.Path == path || len(r.URL.Path) > len(path) && r.URL.Path[:len(path)] == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			authToken := r.Header.Get("Authorization")
			if authToken != "Bearer "+token {
				handlers.WriteJSONError(w, http.StatusUnauthorized, "Não autorizado")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
