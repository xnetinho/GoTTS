package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"tts-api/internal/voice"
)

type TTSHandler struct {
	voiceManager *voice.Manager
}

type SynthesizeRequest struct {
	Text  string `json:"text"`
	Voice string `json:"voice"`
}

func NewTTSHandler(vm *voice.Manager) *TTSHandler {
	return &TTSHandler{voiceManager: vm}
}

func (h *TTSHandler) Synthesize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req SynthesizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao ler requisição", http.StatusBadRequest)
		return
	}

	// Validações adicionais
	if req.Text == "" {
		http.Error(w, "Texto não pode estar vazio", http.StatusBadRequest)
		return
	}

	if req.Voice == "" {
		voices := h.voiceManager.ListVoices()
		http.Error(w, fmt.Sprintf("Voz não especificada. Vozes disponíveis: %v", voices), http.StatusBadRequest)
		return
	}

	audio, err := h.voiceManager.Synthesize(req.Text, req.Voice)
	if err != nil {
		voices := h.voiceManager.ListVoices()
		if strings.Contains(err.Error(), "não encontrada") {
			http.Error(w, fmt.Sprintf("%v. Vozes disponíveis: %v", err, voices), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(audio)))
	w.Write(audio)
}

func (h *TTSHandler) ListVoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	voices := h.voiceManager.ListVoices()
	json.NewEncoder(w).Encode(map[string][]string{"voices": voices})
}
