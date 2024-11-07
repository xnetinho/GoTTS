package config

import (
	"os"
	"strings"
)

type Config struct {
	Port       string
	AuthToken  string
	VoiceFiles []string
	VoicesDir  string
}

func Load() *Config {
	return &Config{
		Port:       getEnvOrDefault("PORT", "8080"),
		AuthToken:  getEnvOrDefault("AUTH_TOKEN", "default-token"),
		VoiceFiles: strings.Split(getEnvOrDefault("VOICE_FILES", "en_US-jenny-medium.onnx"), ","),
		VoicesDir:  getEnvOrDefault("VOICES_DIR", "./voices"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
