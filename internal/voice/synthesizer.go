package voice

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Synthesize(voiceDir string, text string) ([]byte, error) {
	// Garantir que o texto termine com pontuação
	if len(text) > 0 && !strings.ContainsAny(text[len(text)-1:], ".!?") {
		text = text + "."
	}

	// Caminhos para os arquivos do modelo e configuração
	modelPath := ""
	configPath := ""

	// Procurar pelos arquivos .onnx e .onnx.json no diretório da voz
	files, err := filepath.Glob(filepath.Join(voiceDir, "*.onnx"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("nenhum arquivo .onnx encontrado na voz %s", voiceDir)
	}
	modelPath = files[0]
	configPath = modelPath + ".json"

	// Verificar se o arquivo de configuração existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo de configuração não encontrado: %s", configPath)
	} else if err != nil {
		return nil, fmt.Errorf("erro ao verificar o arquivo de configuração: %v", err)
	}

	// Executar o binário do piper
	cmd := exec.Command("piper",
		"--model", modelPath,
		"--config", configPath,
		"--output_file", "-",
	)
	cmd.Stdin = strings.NewReader(text)

	// Capturar a saída
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Executar o comando
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("erro na síntese: %v: %s", err, stderr.String())
	}

	return out.Bytes(), nil
}
