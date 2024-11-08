package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	// Validação do tamanho do texto
	if len(req.Text) > h.voiceManager.Config.MaxTexto {
		mensagem := map[string]interface{}{
			"erro":         "O texto enviado excede o limite estabelecido",
			"limite":       h.voiceManager.Config.MaxTexto,
			"tamanhoTexto": len(req.Text),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(mensagem)
		return
	}

	if req.Voice == "" {
		voices := h.voiceManager.ListVoices()
		http.Error(w, fmt.Sprintf("Voz não especificada. Vozes disponíveis: %v", voices), http.StatusBadRequest)
		return
	}

	// Obter o formato solicitado
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "base64" // Padrão é base64
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

	if format == "base64" {
		// Codificar o áudio em base64
		encodedAudio := base64.StdEncoding.EncodeToString(audio)
		w.Header().Set("Content-Type", "application/base64")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(encodedAudio))
	} else if format == "binary" {
		// Enviar o áudio em binário
		w.Header().Set("Content-Type", "audio/wav")
		w.Header().Set("Content-Length", strconv.Itoa(len(audio)))
		w.WriteHeader(http.StatusOK)
		w.Write(audio)
	} else {
		// Formato não suportado
		http.Error(w, "Formato inválido. Use 'binary' ou 'base64'.", http.StatusBadRequest)
	}
}

func (h *TTSHandler) ListVoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	voices := h.voiceManager.ListVoices()
	json.NewEncoder(w).Encode(map[string][]string{"voices": voices})
}
