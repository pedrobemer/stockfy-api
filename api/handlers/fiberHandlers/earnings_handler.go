package fiberHandlers

import (
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"
	"time"

	"github.com/gofiber/fiber/v2"
)

type EarningsApi struct {
	ApplicationLogic usecases.Applications
	ApiLogic         logicApi.UseCases
}

func (earnings *EarningsApi) CreateEarnings(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	var earningsInsert presenter.EarningsBody
	if err := c.BodyParser(&earningsInsert); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	httpStatusCode, earningsCreated, err := earnings.ApiLogic.ApiCreateEarnings(
		earningsInsert.Symbol, earningsInsert.Currency, earningsInsert.EarningType,
		earningsInsert.Date, earningsInsert.Amount, userId.String())

	if httpStatusCode == 400 {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if httpStatusCode == 404 {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiAssetSymbolUser.Error(),
			"code":    404,
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

	earningsApiReturn := presenter.ConvertEarningToApiReturn(earningsCreated.Id,
		earningsInsert.EarningType, earningsCreated.Earning, earningsCreated.Currency,
		earningsCreated.Date, earningsCreated.Asset.Id, earningsCreated.Asset.Symbol)

	err = c.JSON(&fiber.Map{
		"success": true,
		"earning": earningsApiReturn,
		"message": "Earning registered successfully",
	})

	return err

}

func (earnings *EarningsApi) GetEarningsFromAssetUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	httpStatusCode, earningsReturned, err := earnings.ApiLogic.
		ApiGetEarningsFromAssetUser(c.Query("symbol"), userId.String(),
			c.Query("orderBy"), c.Query("limit"), c.Query("offset"))

	if httpStatusCode == 400 {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if httpStatusCode == 404 {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"code":    404,
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

	earningsApiReturn := presenter.ConvertArrayEarningToApiReturn(earningsReturned)

	err = c.JSON(&fiber.Map{
		"success": true,
		"earning": earningsApiReturn,
		"message": "Earnings returned successfully",
	})

	return err
}

func (earnings *EarningsApi) DeleteEarningFromUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	earningId, err := earnings.ApplicationLogic.EarningsApp.DeleteEarningsFromUser(
		c.Params("id"), userId.String())
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	earningApiReturn := presenter.ConvertEarningToApiReturn(*earningId, "", 0,
		"", time.Time{}, "", "")

	err = c.JSON(&fiber.Map{
		"success": true,
		"earning": earningApiReturn,
		"message": "Earning deleted successfully",
	})

	return err
}

func (earnings *EarningsApi) UpdateEarningFromUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	var earningsUpdate presenter.EarningsBody
	if err := c.BodyParser(&earningsUpdate); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	httpStatusCode, updatedEarnings, err := earnings.ApiLogic.
		ApiUpdateEarningsFromUser(c.Params("id"), earningsUpdate.Amount,
			earningsUpdate.EarningType, earningsUpdate.Date, userId.String())

	if httpStatusCode == 400 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if httpStatusCode == 404 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"code":    404,
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

	earningsApiReturn := presenter.ConvertEarningToApiReturn(updatedEarnings.Id,
		updatedEarnings.Type, updatedEarnings.Earning, updatedEarnings.Currency,
		updatedEarnings.Date, updatedEarnings.Asset.Id, updatedEarnings.Asset.Symbol)

	err = c.JSON(&fiber.Map{
		"success": true,
		"earning": earningsApiReturn,
		"message": "Earning updated successfully",
	})

	return err
}
