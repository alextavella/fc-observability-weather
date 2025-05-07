package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	// Shutdown gracefully on interrupt signal
	shuCh := make(chan os.Signal, 1)
	go func() {
		<-shuCh
		log.Println("Recebido sinal de desligamento, encerrando...")
		if err := otelShutdown(); err != nil {
			log.Fatalf("Erro ao desligar OpenTelemetry: %v", err)
		}
		log.Println("OpenTelemetry desligado com sucesso")
		os.Exit(0)
	}()

	// HTTP Server
	app := fiber.New()
	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendString("OK")
	})
	log.Println("Servidor ouvindo em :8080")
	log.Fatal(app.Listen(":8080"))
}
