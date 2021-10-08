package fiberHandlers

import (
	"fmt"
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"

	"github.com/gofiber/fiber/v2"
)

type EarningsApi struct {
	ApplicationLogic usecases.Applications
	ApiLogic         logicApi.Application
}

func (earnings *EarningsApi) CreateEarnings(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	var earningsInsert presenter.EarningsBody
	if err := c.BodyParser(&earningsInsert); err != nil {
		fmt.Println(err)
	}

	httpStatusCode, earningsCreated, err := earnings.ApiLogic.ApiCreateEarnings(
		earningsInsert.Symbol, earningsInsert.Currency, earningsInsert.EarningType,
		earningsInsert.Date, earningsInsert.Amount, userId.String())
	if err != nil {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	earningsApiReturn := presenter.ConvertEarningToApiReturn(earningsCreated.Id,
		earningsInsert.EarningType, earningsCreated.Earning, earningsCreated.Currency,
		&earningsCreated.Date, earningsCreated.Asset.Id, earningsCreated.Asset.Symbol)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": earningsApiReturn,
		"message": "Earning registered successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err

}

func (earnings *EarningsApi) GetEarningsFromAssetUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	httpStatusCode, earningsReturned, err := earnings.ApiLogic.
		ApiGetEarningsFromAssetUser(c.Query("symbol"), userId.String())
	if err != nil {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	earningsApiReturn := presenter.ConvertArrayEarningToApiReturn(earningsReturned)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": earningsApiReturn,
		"message": "Earnings returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

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
			"message": entity.ErrInvalidEarningId.Error(),
		})
	}

	earningApiReturn := presenter.ConvertEarningToApiReturn(*earningId, "", 0,
		"", nil, "", "")

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": earningApiReturn,
		"message": "Earning deleted successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err
}

func (earnings *EarningsApi) UpdateEarningFromUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	var earningsUpdate presenter.EarningsBody
	if err := c.BodyParser(&earningsUpdate); err != nil {
		fmt.Println(err)
	}
	// fmt.Println(earningsUpdate)

	httpStatusCode, updatedEarnings, err := earnings.ApiLogic.
		ApiUpdateEarningsFromUser(c.Params("id"), earningsUpdate.Amount,
			earningsUpdate.EarningType, earningsUpdate.Date, userId.String())
	if err != nil {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	earningsApiReturn := presenter.ConvertEarningToApiReturn(updatedEarnings.Id,
		updatedEarnings.Type, updatedEarnings.Earning, updatedEarnings.Currency,
		&updatedEarnings.Date, updatedEarnings.Asset.Id, updatedEarnings.Asset.Symbol)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": earningsApiReturn,
		"message": "Earning updated successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err
}
