package handlers

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		writeJSONError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	var req SynthesizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Erro ao ler requisição")
		return
	}

	// Validações adicionais
	if req.Text == "" {
		writeJSONError(w, http.StatusBadRequest, "Texto não pode estar vazio")
		return
	}

	// Validação do tamanho do texto
	if len(req.Text) > h.voiceManager.Config.MaxTexto {
		mensagem := map[string]interface{}{
			"erro":         "O texto enviado excede o limite estabelecido",
			"limite":       h.voiceManager.Config.MaxTexto,
			"tamanhoTexto": len(req.Text),
		}
		writeJSONResponse(w, http.StatusBadRequest, mensagem)
		return
	}

	if req.Voice == "" {
		voices := h.voiceManager.ListVoices()
		mensagem := map[string]interface{}{
			"erro":             "Voz não especificada",
			"vozesDisponiveis": voices,
		}
		writeJSONResponse(w, http.StatusBadRequest, mensagem)
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
		mensagem := map[string]interface{}{
			"erro":             err.Error(),
			"vozesDisponiveis": voices,
		}
		writeJSONResponse(w, http.StatusBadRequest, mensagem)
		return
	}

	// Calcular a duração do áudio
	duration, err := calculateWavDuration(audio)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Erro ao calcular a duração do áudio: %v", err))
		return
	}

	if format == "binary" {
		// Retornar o áudio binário diretamente
		w.Header().Set("Content-Type", "audio/wav")
		w.Header().Set("Content-Length", strconv.Itoa(len(audio)))
		w.Header().Set("X-Duration-Seconds", fmt.Sprintf("%.2f", duration))
		w.WriteHeader(http.StatusOK)
		w.Write(audio)
	} else {
		// Codificar o áudio em base64 e retornar em JSON
		encodedAudio := base64.StdEncoding.EncodeToString(audio)
		response := map[string]interface{}{
			"duration": duration,
			"voice":    req.Voice,
			"text":     req.Text,
			"audio":    encodedAudio,
		}
		writeJSONResponse(w, http.StatusOK, response)
	}
}

func (h *TTSHandler) ListVoices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "Método não permitido")
		return
	}

	voices := h.voiceManager.ListVoices()
	writeJSONResponse(w, http.StatusOK, map[string][]string{"voices": voices})
}

// Função para calcular a duração do áudio em segundos
func calculateWavDuration(data []byte) (float64, error) {
	if len(data) < 44 {
		return 0, fmt.Errorf("dados insuficientes para um arquivo WAV válido")
	}

	// Ler o número de canais (2 bytes a partir do byte 22)
	numChannels := binary.LittleEndian.Uint16(data[22:24])

	// Ler a taxa de amostragem (4 bytes a partir do byte 24)
	sampleRate := binary.LittleEndian.Uint32(data[24:28])

	// Ler bits por amostra (2 bytes a partir do byte 34)
	bitsPerSample := binary.LittleEndian.Uint16(data[34:36])

	// Calcular o tamanho dos dados de áudio (tamanho total menos o header de 44 bytes)
	dataSize := len(data) - 44

	// Calcular o número total de amostras
	bytesPerSample := bitsPerSample / 8
	totalSamples := uint32(dataSize) / uint32(bytesPerSample*numChannels)

	// Calcular a duração
	duration := float64(totalSamples) / float64(sampleRate)

	return duration, nil
}

// Função auxiliar para escrever respostas JSON
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Função auxiliar para escrever erros em JSON
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	writeJSONResponse(w, statusCode, map[string]string{"erro": message})
}
