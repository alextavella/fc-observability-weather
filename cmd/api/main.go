package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alextavella/fc-observability-weather/internal/config"
	"github.com/alextavella/fc-observability-weather/internal/handler"
	"github.com/alextavella/fc-observability-weather/internal/lib"
	"github.com/alextavella/fc-observability-weather/internal/service"
	"github.com/alextavella/fc-observability-weather/pkg/otel"
)

func main() {
	cfg, err := config.NewConfig("build", ".env")
	if err != nil {
		log.Fatalf("Erro ao carregar configuração do servidor A: %v", err)
	}

	// OpenTelemetry
	otelShutdown, err := otel.SetupOtel(context.Background(), cfg.OTEL_HOST, "weather-app")
	if err != nil {
		log.Fatalf("Erro ao configurar OpenTelemetry: %v", err)
	}

	// HTTP Servers
	appA := lib.NewServer()
	appB := lib.NewServer()

	// Services for Server A
	httpClient := lib.NewHttpClient()
	appService := service.NewAppService(httpClient, cfg.APP_B_HOST)

	// Handlers for Server A
	zipcodeHandler := handler.NewZipcodeHandler(appService)
	zipcodeHandler.RegisterRoutes(appA)

	// Services for Server B
	addressService := service.NewViaCepService()
	weatherService := service.NewWeatherService(cfg.WEATHER_API_KEY)

	// Handlers for Server B
	weatherHandler := handler.NewWeatherHandler(addressService, weatherService)
	weatherHandler.RegisterRoutes(appB)

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Println(fmt.Printf("Servidor A ouvindo em :%d", cfg.APP_A_PORT))
		if err := appA.Listen(fmt.Sprintf(":%d", cfg.APP_A_PORT)); err != nil {
			log.Fatalf("Erro ao iniciar servidor A: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		log.Println(fmt.Printf("Servidor B ouvindo em :%d", cfg.APP_B_PORT))
		if err := appB.Listen(fmt.Sprintf(":%d", cfg.APP_B_PORT)); err != nil {
			log.Fatalf("Erro ao iniciar servidor B: %v", err)
		}
	}()

	go func() {
		<-sigCh
		log.Println("Received shutdown signal, shutting down...")
		if err := otelShutdown(); err != nil {
			log.Fatalf("Error shutting down OpenTelemetry: %v", err)
		}
		log.Println("OpenTelemetry shut down successfully")
		if err := appA.Shutdown(); err != nil {
			log.Fatalf("Error shutting down server A: %v", err)
		}
		if err := appB.Shutdown(); err != nil {
			log.Fatalf("Error shutting down server B: %v", err)
		}
		log.Println("Servers shut down gracefully")
	}()

	wg.Wait()
}
