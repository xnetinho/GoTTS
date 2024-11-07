package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	voicesJsonURL   = "https://huggingface.co/rhasspy/piper-voices/raw/main/voices.json"
	baseDownloadURL = "https://huggingface.co/rhasspy/piper-voices/resolve/main/"
)

type LanguageInfo struct {
	Code           string `json:"code"`
	Family         string `json:"family"`
	Region         string `json:"region"`
	NameNative     string `json:"name_native"`
	NameEnglish    string `json:"name_english"`
	CountryEnglish string `json:"country_english"`
}

type FileInfo struct {
	SizeBytes int64  `json:"size_bytes"`
	MD5Digest string `json:"md5_digest"`
}

type VoiceInfo struct {
	Key      string              `json:"key"`
	Name     string              `json:"name"`
	Language LanguageInfo        `json:"language"`
	Quality  string              `json:"quality"`
	Files    map[string]FileInfo `json:"files"`
}

type VoicesManifest map[string]VoiceInfo

func DownloadVoices(voicesDir string, requestedVoices []string) error {
	if err := os.MkdirAll(voicesDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório de vozes: %v", err)
	}

	manifest, err := fetchVoicesManifest()
	if err != nil {
		return err
	}

	for _, voiceName := range requestedVoices {
		voiceName = strings.TrimSpace(voiceName)
		found := false

		for _, voice := range manifest {
			if voice.Name == voiceName {
				found = true
				if err := downloadVoiceFiles(voice, voicesDir); err != nil {
					log.Printf("Erro ao baixar voz %s: %v", voiceName, err)
					continue
				}
				break
			}
		}

		if !found {
			log.Printf("Aviso: voz %s não encontrada no manifesto", voiceName)
		}
	}

	return nil
}

func fetchVoicesManifest() (VoicesManifest, error) {
	resp, err := http.Get(voicesJsonURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar manifesto de vozes: %v", err)
	}
	defer resp.Body.Close()

	var manifest VoicesManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("erro ao decodificar manifesto: %v", err)
	}

	return manifest, nil
}

func downloadVoiceFiles(voice VoiceInfo, voicesDir string) error {
	voiceDir := filepath.Join(voicesDir, voice.Name)
	if err := os.MkdirAll(voiceDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório para a voz %s: %v", voice.Name, err)
	}

	for filePath := range voice.Files {
		url := baseDownloadURL + filePath
		filename := filepath.Base(filePath)
		targetPath := filepath.Join(voiceDir, filename)

		if err := downloadFile(url, targetPath); err != nil {
			return err
		}
		log.Printf("Arquivo baixado com sucesso: %s", filename)
	}
	return nil
}

func downloadFile(url, targetPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
