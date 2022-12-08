package fiberHandlers

import (
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type EventsApi struct {
	ApplicationLogic   usecases.Applications
	ExternalInterfaces externalapi.ThirdPartyInterfaces
	LogicApi           logicApi.UseCases
}

func (events *EventsApi) CreateEventOrder(c *fiber.Ctx) error {

	var eventInserted presenter.EventBody
	if err := c.BodyParser(&eventInserted); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	httpStatusCode, eventsCreated, err := events.LogicApi.ApiCreateEvent(
		eventInserted.Symbol, eventInserted.SymbolDemerger,
		eventInserted.EventType, eventInserted.EventRate, eventInserted.Price,
		eventInserted.Currency, eventInserted.Date, userId.String())

	if httpStatusCode == 400 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if httpStatusCode == 500 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
		})
	}

	eventsApiReturn := presenter.ConvertOrderToApiReturn(eventsCreated)

	err = c.JSON(&fiber.Map{
		"success": true,
		"orders":  eventsApiReturn,
		"message": "Event registered successfully",
	})

	return err
}

func (events *EventsApi) UpdateEventOrder(c *fiber.Ctx) error {
	return nil
}

func (events *EventsApi) DeleteEventOrder(c *fiber.Ctx) error {
	return nil
}

func (events *EventsApi) GetEventOrder(c *fiber.Ctx) error {
	return nil
}
