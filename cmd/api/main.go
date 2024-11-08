package main

import (
	"log"
	"net/http"
	"tts-api/internal/config"
	"tts-api/internal/handlers"
	"tts-api/internal/middleware"
	"tts-api/internal/voice"
	"tts-api/internal/voice/downloader"
)

func main() {
	cfg := config.Load()

	// Download das vozes solicitadas
	if err := downloader.DownloadVoices(cfg.VoicesDir, cfg.Voices); err != nil {
		log.Printf("Aviso: erro no download das vozes: %v", err)
	}

	// Inicializa o gerenciador com as vozes disponíveis
	voiceManager, err := voice.NewManager(cfg)
	if err != nil {
		log.Fatalf("Falha ao inicializar gerenciador de vozes: %v", err)
	}

	// Lista as vozes disponíveis
	voices := voiceManager.ListVoices()
	log.Printf("Vozes disponíveis: %v", voices)

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
