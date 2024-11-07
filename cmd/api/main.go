package main

import (
	"log"
	"net/http"
	"tts-api/internal/config"
	"tts-api/internal/handlers"
	"tts-api/internal/middleware"
	"tts-api/internal/voice"
)

func main() {
	cfg := config.Load()

	voiceManager, err := voice.NewManager(cfg.VoicesDir, cfg.VoiceFiles)
	if err != nil {
		log.Fatalf("Falha ao inicializar gerenciador de vozes: %v", err)
	}
	defer voiceManager.Close() // Importante: adicionar esta linha

	ttsHandler := handlers.NewTTSHandler(voiceManager)

	mux := http.NewServeMux()
	mux.HandleFunc("/synthesize", ttsHandler.Synthesize)
	mux.HandleFunc("/voices", ttsHandler.ListVoices)

	handler := middleware.AuthMiddleware(cfg.AuthToken)(mux)

	log.Printf("Servidor iniciando na porta %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
