package voice

import (
	"fmt"
	"os"
	"strings"

	"github.com/amitybell/piper"
	asset "github.com/amitybell/piper-asset"
)

func Synthesize(voiceDir string, text string) ([]byte, error) {
	// Garantir que o texto termine com pontuação
	if len(text) > 0 && !strings.ContainsAny(text[len(text)-1:], ".!?") {
		text = text + "."
	}

	// Criar um asset.Asset personalizado usando o diretório da voz
	voiceAsset := asset.Asset{
		Name: "custom-voice",
		FS:   os.DirFS(voiceDir),
	}

	// Especificar o dataDir (pode ser vazio ou um caminho específico)
	dataDir := ""

	// Criar nova instância do TTS com os argumentos corretos
	tts, err := piper.New(dataDir, voiceAsset)
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
