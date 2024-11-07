package voice

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/amitybell/piper"
	asset "github.com/amitybell/piper-asset"
)

func Synthesize(modelPath string, text string) ([]byte, error) {
	// Garantir que o texto termine com pontuação
	if len(text) > 0 && !strings.ContainsAny(text[len(text)-1:], ".!?") {
		text = text + "."
	}

	// Criar o asset a partir do arquivo
	modelAsset := asset.NewFile(modelPath)

	// Obter diretório absoluto para instalação
	dataDir, err := filepath.Abs(filepath.Dir(modelPath))
	if err != nil {
		return nil, fmt.Errorf("erro ao obter path absoluto: %v", err)
	}

	// Criar nova instância do TTS usando o sistema de instalação do piper
	tts, err := piper.New(dataDir, modelAsset)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar TTS: %v", err)
	}

	// Sintetizar o texto
	audio, err := tts.Synthesize(text)
	if err != nil {
		return nil, fmt.Errorf("erro na síntese: %v", err)
	}

	return audio, nil
}
