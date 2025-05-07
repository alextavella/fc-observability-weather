package handler

import (
	"fmt"

	"github.com/alextavella/fc-observability-weather/internal/config"
	"github.com/alextavella/fc-observability-weather/internal/service"
	"github.com/alextavella/fc-observability-weather/pkg/otel"
	"github.com/gofiber/fiber/v3"
)

type zipcodeHandler struct {
	weatherService service.IWeatherService
}

func NewZipcodeHandler(cfg *config.Config) *zipcodeHandler {
	return &zipcodeHandler{
		weatherService: service.NewWeatherService(cfg.APP_B_HOST),
	}
}
func (h *zipcodeHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/zipcode", h.HandleZipCodeRequest)
}

type zipcodeInput struct {
	Zipcode string `json:"cep"`
}

func (i *zipcodeInput) Validate() error {
	if len(i.Zipcode) != 8 {
		return fmt.Errorf("invalid zipcode")
	}
	return nil
}

func (h *zipcodeHandler) HandleZipCodeRequest(c fiber.Ctx) error {
	tracer := otel.GetTracer()
	ctx, span := tracer.Start(c.Context(), "zipcode-request")
	defer span.End()

	input := &zipcodeInput{}
	if err := c.Bind().Body(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate the input
	err := input.Validate()
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Call the weather service
	weatherResp, err := h.weatherService.GetWeatherByZipCode(ctx, input.Zipcode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fmt.Println("Weather response:", weatherResp)

	return c.JSON(fiber.Map{
		"cep": input.Zipcode,
	})
}
