package handlers

import (
	"fmt"
	"stockfyApi/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func PostOrder(c *fiber.Ctx) error {

	var orderInsert database.OrderBodyPost
	if err := c.BodyParser(&orderInsert); err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderInsert)

	var brokerageId string
	var orderReturn database.OrderApiReturn
	var err error
	var assetExist bool
	var Ids DatabaseId

	var condAssetExist = "symbol='" + orderInsert.Symbol + "'"
	assetExist = database.VerifyRowExistence(*database.DBpool, "asset",
		condAssetExist)

	if !assetExist {
		var apiType string

		if orderInsert.Country == "BR" {
			apiType = "Alpha"
		} else {
			apiType = "Finnhub"
		}
		_, Ids = assetVerification(orderInsert.Symbol, orderInsert.Country,
			apiType)

		if Ids.AssetId == "" {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": "The Symbol " + orderInsert.Symbol +
					" from country " + orderInsert.Country + " do not exist.",
			})
		}
	} else {
		symbolQuery := database.SearchAsset(database.DBpool,
			orderInsert.Symbol, "")
		Ids.AssetId = symbolQuery[0].Id
	}

	brokerageReturn, _ := database.FetchBrokerage(*database.DBpool, "SINGLE",
		orderInsert.Brokerage)
	brokerageId = brokerageReturn[0].Id

	orderReturn = database.CreateOrder(*database.DBpool, orderInsert,
		Ids.AssetId, brokerageId)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"orders":  orderReturn,
		"message": "Order registered successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func DeleteOrder(c *fiber.Ctx) error {
	var err error

	orderId := database.DeleteOrder(*database.DBpool, c.Params("id"))
	if orderId == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "The order " + c.Params("id") +
				" does not exist in your table. Please provide a valid ID.",
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"order":   orderId,
		"message": "Order deleted successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
