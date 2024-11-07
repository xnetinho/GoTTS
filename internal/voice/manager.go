package voice

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nabbl/piper"
)

type Manager struct {
	voices    map[string]*piper.TTS
	voicesDir string
	mu        sync.RWMutex
}

func NewManager(dataDir string, voiceFiles []string) (*Manager, error) {
	// Verificar e criar diretório se não existir
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório de vozes: %v", err)
	}

	m := &Manager{
		voices:    make(map[string]*piper.TTS),
		voicesDir: dataDir,
	}

	for _, voiceFile := range voiceFiles {
		// Limpar e validar nome do arquivo
		voiceFile = strings.TrimSpace(voiceFile)
		if !strings.HasSuffix(strings.ToLower(voiceFile), ".onnx") {
			log.Printf("Aviso: arquivo ignorado %s - extensão inválida", voiceFile)
			continue
		}

		// Construir caminho absoluto
		voicePath := filepath.Clean(filepath.Join(dataDir, voiceFile))

		// Verificar se o arquivo existe
		if _, err := os.Stat(voicePath); err != nil {
			if os.IsNotExist(err) {
				log.Printf("Aviso: arquivo de voz não encontrado: %s", voicePath)
				continue
			}
			return nil, fmt.Errorf("erro ao verificar arquivo %s: %v", voicePath, err)
		}

		// Carregar modelo TTS
		tts, err := piper.New(voicePath)
		if err != nil {
			return nil, fmt.Errorf("falha ao carregar voz %s: %v", voiceFile, err)
		}

		// Usar nome base do arquivo como identificador
		voiceName := filepath.Base(voiceFile)
		m.voices[voiceName] = tts
		log.Printf("Voz carregada com sucesso: %s", voiceName)
	}

	if len(m.voices) == 0 {
		return nil, fmt.Errorf("nenhuma voz válida foi carregada do diretório %s", dataDir)
	}

	return m, nil
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for _, tts := range m.voices {
		if err := tts.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func (m *Manager) Synthesize(text, voice string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tts, exists := m.voices[voice]
	if !exists {
		return nil, fmt.Errorf("voz %s não encontrada", voice)
	}

	return tts.Synthesize(text)
}

func (m *Manager) ListVoices() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	voices := make([]string, 0, len(m.voices))
	for voice := range m.voices {
		voices = append(voices, voice)
	}
	return voices
}

// Adicione este método para obter o diretório de vozes
func (m *Manager) GetVoicesDir() string {
	return m.voicesDir
}
