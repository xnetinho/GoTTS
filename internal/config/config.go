package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port      string
	AuthToken string
	Voices    []string
	VoicesDir string
	MaxTexto  int // Novo campo adicionado
}

func Load() *Config {
	maxTextoStr := getEnvOrDefault("MAX_TEXTO", "100000")
	maxTexto, err := strconv.Atoi(maxTextoStr)
	if err != nil {
		maxTexto = 100000
	}

	return &Config{
		Port:      getEnvOrDefault("PORT", "8080"),
		AuthToken: getEnvOrDefault("AUTH_TOKEN", "default-token"),
		Voices:    strings.Split(getEnvOrDefault("VOICE_FILES", ""), ","),
		VoicesDir: getEnvOrDefault("VOICES_DIR", "./voices"),
		MaxTexto:  maxTexto, // Atribui o valor lido
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
