package main

import (
	"log"
	"net/http"
	"tts-api/internal/config"
	"tts-api/internal/handlers"
	"tts-api/internal/middleware"
	"tts-api/internal/voice"
	"tts-api/internal/voice/downloader"

	_ "tts-api/docs" // Importa o pacote docs gerado pelo swag

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title GoTTS API
// @version 1.0
// @description API para síntese de voz usando Go e Piper

// @contact.name Diomedes Neto
// @contact.url http://www.izap.app
// @contact.email admin@izap.app

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

	// Rotas que não exigem autenticação
	mux.HandleFunc("/healthcheck", handlers.HealthCheckHandler)
	mux.HandleFunc("/api/", httpSwagger.WrapHandler)

	// Rotas que exigem autenticação
	mux.HandleFunc("/synthesize", ttsHandler.Synthesize)
	mux.HandleFunc("/voices", ttsHandler.ListVoices)

	// Aplica o middleware de autenticação nas rotas que exigem
	handler := middleware.AuthMiddleware(cfg.AuthToken)(mux)

	log.Printf("Servidor iniciando na porta %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
