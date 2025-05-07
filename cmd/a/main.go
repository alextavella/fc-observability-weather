package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alextavella/fc-observability-weather/internal/config"
	"github.com/alextavella/fc-observability-weather/internal/handler"
	"github.com/alextavella/fc-observability-weather/pkg/otel"
	"github.com/gofiber/fiber/v3"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig("build", ".env.a")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	//  OpenTelemetry configuration
	otelShutdown, err := otel.SetupOtel(ctx, cfg.OTEL_HOST, cfg.APP_NAME)
	if err != nil {
		log.Fatalf("Erro ao configurar OpenTelemetry: %v", err)
	}

	// HTTP Server
	app := fiber.New()

	// Handlers
	zipcodeHandler := handler.NewZipcodeHandler(cfg)
	zipcodeHandler.RegisterRoutes(app)

	// Shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Received shutdown signal, shutting down...")
		if err := otelShutdown(); err != nil {
			log.Fatalf("Error shutting down OpenTelemetry: %v", err)
		}
		log.Println("OpenTelemetry shut down successfully")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error shutting down: %v", err)
		}
		log.Println("Server shut down gracefully")
	}()

	// Start server
	log.Println(fmt.Printf("Servidor ouvindo em :%d", cfg.PORT))
	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.PORT)))
}
