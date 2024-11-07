package voice

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/amitybell/piper"
	asset "github.com/amitybell/piper-asset"
)

type Manager struct {
	voices    map[string]*piper.TTS
	voicesDir string
	mu        sync.RWMutex
}

func NewManager(dataDir string, voiceFiles []string) (*Manager, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório de vozes: %v", err)
	}

	m := &Manager{
		voices:    make(map[string]*piper.TTS),
		voicesDir: dataDir,
	}

	for _, voiceFile := range voiceFiles {
		voiceFile = strings.TrimSpace(voiceFile)
		if !strings.HasSuffix(strings.ToLower(voiceFile), ".onnx") {
			log.Printf("Aviso: arquivo ignorado %s - extensão inválida", voiceFile)
			continue
		}

		voicePath := filepath.Clean(filepath.Join(dataDir, voiceFile))

		if _, err := os.Stat(voicePath); err != nil {
			if os.IsNotExist(err) {
				log.Printf("Aviso: arquivo de voz não encontrado: %s", voicePath)
				continue
			}
			return nil, fmt.Errorf("erro ao verificar arquivo %s: %v", voicePath, err)
		}

		modelAsset := asset.NewFile(voicePath)
		tts, err := piper.New(voicePath, modelAsset)
		if err != nil {
			return nil, fmt.Errorf("falha ao carregar voz %s: %v", voiceFile, err)
		}

		voiceName := filepath.Base(voiceFile)
		m.voices[voiceName] = tts
		log.Printf("Voz carregada com sucesso: %s", voiceName)
	}

	return m, nil
}

func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, v := range m.voices {
		v.Close()
	}
}

func (m *Manager) Synthesize(text, voice string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tts, exists := m.voices[voice]
	if !exists {
		return nil, fmt.Errorf("voz %s não encontrada", voice)
	}

	if text == "" {
		return nil, fmt.Errorf("texto não pode estar vazio")
	}

	// Garantir que o texto termine com pontuação
	if !strings.ContainsAny(text[len(text)-1:], ".!?") {
		text = text + "."
	}

	audio, err := tts.Synthesize(text)
	if err != nil {
		return nil, fmt.Errorf("erro na síntese: %v", err)
	}

	return audio, nil
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

func (m *Manager) GetVoicesDir() string {
	return m.voicesDir
}
