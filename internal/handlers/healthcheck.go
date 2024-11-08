package handlers

import (
	"net/http"
)

// HealthCheckHandler lida com a rota /healthcheck
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSONError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	response := map[string]string{
		"status": "OK",
	}
	WriteJSONResponse(w, http.StatusOK, response)
}
