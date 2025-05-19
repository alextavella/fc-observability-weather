package handler

import (
	"github.com/alextavella/fc-observability-weather/internal/domain/exception"
	"github.com/alextavella/fc-observability-weather/internal/service"
	"github.com/alextavella/fc-observability-weather/pkg/otel"
	"github.com/gofiber/fiber/v3"
)

type zipcodeHandler struct {
	appService service.IAppService
}

func NewZipcodeHandler(appService service.IAppService) *zipcodeHandler {
	return &zipcodeHandler{
		appService: appService,
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
		return exception.ErrInvalidZipcode
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
			"error": exception.ErrParseBody.Error(),
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
	weatherResp, err := h.appService.GetWeather(ctx, input.Zipcode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": exception.ErrCanNotFindZipcode.Error(),
		})
	}

	return c.JSON(weatherResp)
}
