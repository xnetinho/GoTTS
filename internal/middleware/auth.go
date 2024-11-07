package middleware

import "net/http"

func AuthMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken := r.Header.Get("Authorization")
			if authToken != "Bearer "+token {
				http.Error(w, "Não autorizado", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
