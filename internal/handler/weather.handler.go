package handler

import (
	"github.com/alextavella/fc-observability-weather/internal/service"
	"github.com/alextavella/fc-observability-weather/pkg/otel"
	"github.com/gofiber/fiber/v3"
)

type weatherHandler struct {
	addressService service.IAddressService
}

func NewWeatherHandler() *weatherHandler {
	return &weatherHandler{
		addressService: service.NewViaCepService(),
	}
}
func (h *weatherHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/weather/:zipcode/zipcode", h.HandleWeatherRequest)
}

type weatherInput struct {
	Zipcode string
}

func (h *weatherHandler) HandleWeatherRequest(c fiber.Ctx) error {
	tracer := otel.GetTracer()
	ctx, span := tracer.Start(c.Context(), "weather-request")
	defer span.End()

	input := new(weatherInput)
	input.Zipcode = c.Params("zipcode")

	viacepResult, err := h.addressService.GetAddressByZipcode(ctx, input.Zipcode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"neighborhood": viacepResult.Bairro,
	})
}
