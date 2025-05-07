package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	otel "github.com/alextavella/go-opentelemetry/pkg/otel"
	"github.com/gofiber/fiber/v3"
)

func main() {
	ctx := context.Background()

	name := os.Getenv("APP_NAME")
	fmt.Println("APP_NAME:", name)
	if name == "" {
		name = "weather-app"
	}

	//  OpenTelemetry configuration
	otelShutdown, err := otel.SetupOtel(ctx, name)
	if err != nil {
		log.Fatalf("Erro ao configurar OpenTelemetry: %v", err)
	}

	// HTTP Server
	app := fiber.New()
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})

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
	log.Println("Servidor ouvindo em :8080")
	log.Fatal(app.Listen(":8080"))
}
