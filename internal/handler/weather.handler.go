package handler

import (
	"github.com/alextavella/fc-observability-weather/internal/service"
	"github.com/alextavella/fc-observability-weather/internal/util"
	"github.com/alextavella/fc-observability-weather/pkg/otel"
	"github.com/gofiber/fiber/v3"
	"go.opentelemetry.io/otel/propagation"
)

type weatherHandler struct {
	addressService service.IAddressService
	weatherService service.IWeatherService
}

func NewWeatherHandler(addressService service.IAddressService, weatherService service.IWeatherService) *weatherHandler {
	return &weatherHandler{
		addressService: addressService,
		weatherService: weatherService,
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

	// Extraindo o contexto de rastreamento dos cabeçalhos
	carrier := propagation.HeaderCarrier(c.GetReqHeaders())
	ctx := propagation.TraceContext{}.Extract(c.Context(), carrier)

	// Iniciando o span com o contexto extraído
	ctx, span := tracer.Start(ctx, "weather-request")
	defer span.End()

	input := new(weatherInput)
	input.Zipcode = c.Params("zipcode")

	viacepResult, err := h.addressService.GetAddressByZipcode(ctx, input.Zipcode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	zipcode := viacepResult.Cep
	weatherResult, err := h.weatherService.GetWeatherByZipCode(ctx, zipcode)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_, spanTemp := tracer.Start(ctx, "convert-temperatures")
	defer spanTemp.End()

	tempC := weatherResult.Current.TempC
	tempF := util.ConvertToFahrenheit(tempC)
	tempK := util.ConvertToKelvin(tempC)

	return c.JSON(fiber.Map{
		"zipcode":    zipcode,
		"celsius":    tempC,
		"fahrenheit": tempF,
		"kelvin":     tempK,
	})
}
