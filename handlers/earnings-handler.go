package handlers

import (
	"fmt"
	"stockfyApi/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type EarningsApi struct {
	Db database.PgxIface
}

func (earnings *EarningsApi) PostEarnings(c *fiber.Ctx) error {
	var err error

	validEarningTypes := map[string]bool{"Dividendos": true, "JCP": true,
		"Rendimentos": true}

	var earningsInsert database.EarningsBodyPost
	if err := c.BodyParser(&earningsInsert); err != nil {
		fmt.Println(err)
	}
	fmt.Println(earningsInsert)

	if earningsInsert.Symbol == "" || earningsInsert.Currency == "" ||
		earningsInsert.EarningType == "" || earningsInsert.Date == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "There is an empty field in the JSON request.",
		})
	}

	if earningsInsert.Amount <= 0 {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "The earning must be higher than 0. The request to save" +
				"it has an earning of " +
				strconv.FormatFloat(earningsInsert.Amount, 'f', -1, 64),
		})

	}

	if !validEarningTypes[earningsInsert.EarningType] {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "The EarningType must be Dividendos, JCP or Rendimentos." +
				"The EarningType sent was " + earningsInsert.EarningType,
		})
	}

	assetInfo, _ := database.SearchAsset(database.DBpool, earningsInsert.Symbol,
		"")
	if assetInfo[0].Id == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "The symbol " + earningsInsert.Symbol +
				" is not registered in the Asset table. Please register it before.",
		})
	}

	earningRow := database.CreateEarningRow(database.DBpool, earningsInsert,
		assetInfo[0].Id)

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
