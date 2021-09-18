package handlers

import (
	"fmt"
	"reflect"
	"stockfyApi/convertVariables"
	"stockfyApi/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type OrderApi struct {
	Db database.PgxIface
}

func (order *OrderApi) PostOrderFromUser(c *fiber.Ctx) error {

	asset := AssetApi{Db: order.Db}

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

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	message := orderVerification(orderInsert)
	if message != "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": message,
		})
	}

	var condAssetExist = "symbol='" + orderInsert.Symbol + "'"
	assetExist = database.VerifyRowExistence(order.Db, "asset",
		condAssetExist)

	if !assetExist {
		var apiType string

		if orderInsert.Country == "BR" {
			apiType = "Alpha"
		} else {
			apiType = "Finnhub"
		}
		_, Ids = asset.assetVerification(orderInsert.Symbol, orderInsert.Country,
			apiType)

		if Ids.AssetId == "" {
			return c.Status(404).JSON(&fiber.Map{
				"success": false,
				"message": "The Symbol " + orderInsert.Symbol +
					" from country " + orderInsert.Country + " do not exist.",
			})
		}
	} else {
		symbolQuery, _ := database.SearchAsset(order.Db, orderInsert.Symbol)
		Ids.AssetId = symbolQuery[0].Id
	}

	assetUser, err := database.SearchAssetUserRelation(order.Db, Ids.AssetId,
		userId.String())
	if len(assetUser) == 0 {
		assetUser, err = database.CreateAssetUserRelation(order.Db, Ids.AssetId,
			userId.String())

		if err != nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
	}

	brokerageReturn, _ := database.FetchBrokerage(order.Db, "SINGLE",
		orderInsert.Brokerage)
	brokerageId = brokerageReturn[0].Id

	orderReturn = database.CreateOrder(order.Db, orderInsert,
		assetUser[0].AssetId, brokerageId, assetUser[0].UserUid)

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

func (order *OrderApi) GetOrdersFromAssetUser(c *fiber.Ctx) error {
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

	asset, err := database.SearchAsset(order.Db, c.Query("symbol"))
	if asset == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "The symbol " + c.Query("symbol") + " does not exist",
		})
	}

	ordersReturn, err := database.SearchOrdersFromAssetUser(order.Db,
		asset[0].Id, userId.String())
	if err != nil {
		if asset == nil {
			return c.Status(500).JSON(&fiber.Map{
				"success": false,
				"message": fmt.Errorf(err.Error()),
			})
		}
	}

	if ordersReturn == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "Your user does not have any order registered " +
				"for the symbol " + c.Query("symbol"),
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"earning": ordersReturn,
		"message": "Orders returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}

func (order *OrderApi) DeleteOrderFromUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	orderId := database.DeleteOrderFromUser(order.Db, c.Params("id"),
		userId.String())
	if orderId == "" {
		return c.Status(404).JSON(&fiber.Map{
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

func (order *OrderApi) UpdateOrderFromUser(c *fiber.Ctx) error {
	var err error

	var orderUpdate database.OrderBodyPost
	if err := c.BodyParser(&orderUpdate); err != nil {
		fmt.Println(err)
	}

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	assetInfo := database.SearchAssetByOrderId(order.Db, orderUpdate.Id)
	fmt.Println(assetInfo)

	if orderUpdate.Id != c.Params("id") {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "The order " + orderUpdate.Id +
				" from the body request is different from the " +
				c.Params("id"),
		})
	}

	orderUpdate.Country = assetInfo[0].AssetType.Country
	message := orderVerification(orderUpdate)
	if message != "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": message,
		})
	}

	orderUpdateReturn := database.UpdateOrderFromUser(order.Db,
		orderUpdate, userId.String())
	if orderUpdateReturn == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "The order Id of " + orderUpdate.Id +
				" does not exist for your user.",
		})
	}

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
