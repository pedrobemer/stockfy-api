package handlers

import (
	"fmt"
	"stockfyApi/convertVariables"
	"stockfyApi/database"
	"strconv"

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

	message := orderVerification(orderInsert)
	if message != "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": message,
		})
	}

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

func UpdateOrder(c *fiber.Ctx) error {
	var err error

	var orderUpdate database.OrderBodyPost
	if err := c.BodyParser(&orderUpdate); err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderUpdate)

	assetInfo := database.SearchAssetByOrderId(*database.DBpool, orderUpdate.Id)
	fmt.Println(assetInfo)

	if orderUpdate.Id != c.Params("id") {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "The order " + orderUpdate.Id +
				" from the body request is different from the " +
				c.Params("id"),
		})
	}

	orderUpdate.Country = assetInfo[0].AssetType.Country
	message := orderVerification(orderUpdate)
	if message != "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": message,
		})
	}

	orderUpdateReturn := database.UpdateOrder(*database.DBpool, orderUpdate)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"order":   orderUpdateReturn,
		"message": "Order updated successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

func orderVerification(orderInfo database.OrderBodyPost) string {
	var message string

	if orderInfo.OrderType != "sell" && orderInfo.OrderType != "buy" {
		message = "The " + orderInfo.OrderType + " does not exist. Please read" +
			" our documentation to understand the possible values for orders."
	}

	if orderInfo.Country == "BR" &&
		(orderInfo.OrderType == "sell" || orderInfo.OrderType == "buy") {
		isIntQuantity := convertVariables.IsIntegral(orderInfo.Quantity)
		if !isIntQuantity {
			message = "The order " + orderInfo.Id +
				" must have a integer quantity value. The " +
				orderInfo.OrderType + " order has a quantity of " +
				strconv.FormatFloat(orderInfo.Quantity, 'f', -1, 64)
		}
	}

	if orderInfo.OrderType == "buy" && orderInfo.Quantity < 0 {
		message = "The order " + orderInfo.Id +
			" must have a positive quantity value. The " +
			orderInfo.OrderType + " order has a quantity of " +
			strconv.FormatFloat(orderInfo.Quantity, 'f', -1, 64)
	} else if orderInfo.OrderType == "sell" && orderInfo.Quantity > 0 {
		message = "The order " + orderInfo.Id +
			" must have a negative quantity value. The " + orderInfo.OrderType +
			" order has a quantity of " +
			strconv.FormatFloat(orderInfo.Quantity, 'f', -1, 64)
	} else if orderInfo.Price < 0 {
		message = "The order " + orderInfo.Id +
			" must have a positive price value. The price sent was " +
			strconv.FormatFloat(orderInfo.Price, 'f', -1, 64)
	}

	return message
}
