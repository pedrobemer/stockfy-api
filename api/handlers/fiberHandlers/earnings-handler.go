package fiberHandlers

import (
	"fmt"
	"reflect"
	"stockfyApi/database"
	"stockfyApi/usecases"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EarningsApi struct {
	ApplicationLogic usecases.Applications
}

func (earnings *EarningsApi) PostEarnings(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	validEarningTypes := map[string]bool{"Dividendos": true, "JCP": true,
		"Rendimentos": true}

	var earningsInsert database.EarningsBodyPost
	if err := c.BodyParser(&earningsInsert); err != nil {
		fmt.Println(err)
	}
	fmt.Println(earningsInsert)

	if earningsInsert.Symbol == "" || earningsInsert.Currency == "" ||
		earningsInsert.EarningType == "" || earningsInsert.Date == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "There is an empty field in the JSON request.",
		})
	}

	if earningsInsert.Amount <= 0 {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "The earning must be higher than 0. The request to save" +
				"it has an earning of " +
				strconv.FormatFloat(earningsInsert.Amount, 'f', -1, 64),
		})

	}

	if !validEarningTypes[earningsInsert.EarningType] {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "The EarningType must be Dividendos, JCP or Rendimentos." +
				"The EarningType sent was " + earningsInsert.EarningType,
		})
	}

	assetInfo, _ := database.SearchAssetByUser(database.DBpool, earningsInsert.Symbol,
		userId.String(), "")
	if assetInfo[0].Id == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "The symbol " + earningsInsert.Symbol +
				" is not registered in the Asset table. Please register it before.",
		})
	}

	earningRow := database.CreateEarningRow(database.DBpool, earningsInsert,
		assetInfo[0].Id, userId.String())

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": earningRow,
		"message": "Earning registered successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func (earnings *EarningsApi) GetEarningsFromAssetUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("symbol") == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "The symbol value in the API query can not be empty. " +
				"Please read our documentation",
		})
	}

	asset, err := database.SearchAsset(earnings.Db, c.Query("symbol"))
	if asset == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "The symbol " + c.Query("symbol") + " does not exist",
		})
	}

	earningsReturn, err := database.SearchEarningFromAssetUser(earnings.Db, asset[0].Id,
		userId.String())
	if err != nil {
		if asset == nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": fmt.Errorf(err.Error()),
			})
		}
	}

	if earningsReturn == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "Your user does not have any earning registered " +
				"for the symbol " + c.Query("symbol"),
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": earningsReturn,
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

	earningId := database.DeleteEarningFromUser(earnings.Db, c.Params("id"),
		userId.String())
	if earningId == "" {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "The earning " + c.Params("id") +
				" does not exist in your table. Please provide a valid ID.",
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": earningId,
		"message": "Order deleted successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

func (earnings *EarningsApi) UpdateEarningFromUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	validEarningTypes := map[string]bool{"Dividendos": true, "JCP": true,
		"Rendimentos": true}

	var earningsUpdate database.EarningsBodyPost
	if err := c.BodyParser(&earningsUpdate); err != nil {
		fmt.Println(err)
	}
	fmt.Println(earningsUpdate)

	if earningsUpdate.EarningType == "" || earningsUpdate.Date == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "There is an empty field in the JSON request.",
		})
	}

	if earningsUpdate.Amount <= 0 {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "The earning must be higher than 0. The request to save" +
				"it has an earning of " +
				strconv.FormatFloat(earningsUpdate.Amount, 'f', -1, 64),
		})

	}

	if !validEarningTypes[earningsUpdate.EarningType] {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "The EarningType must be Dividendos, JCP or Rendimentos." +
				"The EarningType sent was " + earningsUpdate.EarningType,
		})
	}

	earningRow := database.UpdateEarningsFromUser(database.DBpool, earningsUpdate,
		userId.String())

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": earningRow,
		"message": "Earning registered successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
