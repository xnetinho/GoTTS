package handlers

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponse escreve uma resposta JSON com o código de status fornecido
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteJSONError escreve uma mensagem de erro em JSON com o código de status fornecido
func WriteJSONError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSONResponse(w, statusCode, map[string]string{"erro": message})
}
