package voice

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"tts-api/internal/config"
)

type Manager struct {
	voices    map[string]string // mapa de nome -> caminho do arquivo
	voicesDir string
	mu        sync.RWMutex
	Config    *config.Config // Adicionado
}

func NewManager(cfg *config.Config) (*Manager, error) {
	voicesDir := cfg.VoicesDir
	if err := os.MkdirAll(voicesDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório de vozes: %v", err)
	}

	m := &Manager{
		voices:    make(map[string]string),
		voicesDir: voicesDir,
		Config:    cfg, // Atribui a configuração
	}

	entries, err := os.ReadDir(voicesDir)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler diretório de vozes: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			voiceName := entry.Name()
			voicePath := filepath.Join(voicesDir, voiceName)
			m.voices[voiceName] = voicePath
			log.Printf("Voz encontrada: %s", voiceName)
		}
	}

	if len(m.voices) == 0 {
		return nil, fmt.Errorf("nenhuma voz foi encontrada")
	}

	return m, nil
}

func (m *Manager) Synthesize(text, voiceName string) ([]byte, error) {
	m.mu.RLock()
	voiceDir, exists := m.voices[voiceName]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("voz %s não encontrada", voiceName)
	}

	if text == "" {
		return nil, fmt.Errorf("texto não pode estar vazio")
	}

	return Synthesize(voiceDir, text)
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

func (m *Manager) GetVoicePath(voice string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	path, exists := m.voices[voice]
	if !exists {
		return "", fmt.Errorf("voz %s não encontrada", voice)
	}
	return path, nil
}

func (m *Manager) GetVoicesDir() string {
	return m.voicesDir
}

func (m *Manager) Close() {
	// Método mantido vazio para compatibilidade com a interface existente
}
