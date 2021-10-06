package fiberHandlers

import (
	"reflect"
	"stockfyApi/api/presenter"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type OrderApi struct {
	ApplicationLogic   usecases.Applications
	ExternalInterfaces externalapi.ThirdPartyInterfaces
	LogicApi           logicApi.Application
}

func (order *OrderApi) CreateUserOrder(c *fiber.Ctx) error {

	var orderInserted presenter.OrderBody
	if err := c.BodyParser(&orderInserted); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong JSON in the Body",
		})
	}

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	httpStatusCode, orderCreated, err := order.LogicApi.ApiCreateOrder(
		orderInserted.Symbol, orderInserted.Country, orderInserted.OrderType,
		orderInserted.Quantity, orderInserted.Price, orderInserted.Currency,
		orderInserted.Brokerage, orderInserted.Date, userId.String())
	if err != nil {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	orderApiReturn := presenter.ConvertSingleOrderToApiReturn(*orderCreated)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"orders":  orderApiReturn,
		"message": "Order registered successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

// func (order *OrderApi) GetOrdersFromAssetUser(c *fiber.Ctx) error {
// 	var err error

// 	userInfo := c.Context().Value("user")
// 	userId := reflect.ValueOf(userInfo).FieldByName("userID")

// 	if c.Query("symbol") == "" {
// 		return c.Status(400).JSON(&fiber.Map{
// 			"success": false,
// 			"message": "The symbol value in the API query can not be empty. " +
// 				"Please read our documentation",
// 		})
// 	}

// 	asset, err := database.SearchAsset(order.Db, c.Query("symbol"))
// 	if asset == nil {
// 		return c.Status(404).JSON(&fiber.Map{
// 			"success": false,
// 			"message": "The symbol " + c.Query("symbol") + " does not exist",
// 		})
// 	}

// 	ordersReturn, err := database.SearchOrdersFromAssetUser(order.Db,
// 		asset[0].Id, userId.String())
// 	if err != nil {
// 		if asset == nil {
// 			return c.Status(500).JSON(&fiber.Map{
// 				"success": false,
// 				"message": fmt.Errorf(err.Error()),
// 			})
// 		}
// 	}

// 	if ordersReturn == nil {
// 		return c.Status(404).JSON(&fiber.Map{
// 			"success": false,
// 			"message": "Your user does not have any order registered " +
// 				"for the symbol " + c.Query("symbol"),
// 		})
// 	}

// 	if err := c.JSON(&fiber.Map{
// 		"success": true,
// 		"earning": ordersReturn,
// 		"message": "Orders returned successfully",
// 	}); err != nil {
// 		return c.Status(500).JSON(&fiber.Map{
// 			"success": false,
// 			"message": err,
// 		})
// 	}

// 	return err
// }

// func (order *OrderApi) DeleteOrderFromUser(c *fiber.Ctx) error {
// 	var err error

// 	userInfo := c.Context().Value("user")
// 	userId := reflect.ValueOf(userInfo).FieldByName("userID")

// 	orderId := database.DeleteOrderFromUser(order.Db, c.Params("id"),
// 		userId.String())
// 	if orderId == "" {
// 		return c.Status(404).JSON(&fiber.Map{
// 			"success": false,
// 			"message": "The order " + c.Params("id") +
// 				" does not exist in your table. Please provide a valid ID.",
// 		})
// 	}

// 	if err := c.JSON(&fiber.Map{
// 		"success": true,
// 		"order":   orderId,
// 		"message": "Order deleted successfully",
// 	}); err != nil {
// 		return c.Status(500).JSON(&fiber.Map{
// 			"success": false,
// 			"message": err,
// 		})
// 	}

// 	return err
// }

// func (order *OrderApi) UpdateOrderFromUser(c *fiber.Ctx) error {
// 	var err error

// 	var orderUpdate database.OrderBodyPost
// 	if err := c.BodyParser(&orderUpdate); err != nil {
// 		fmt.Println(err)
// 	}

// 	userInfo := c.Context().Value("user")
// 	userId := reflect.ValueOf(userInfo).FieldByName("userID")

// 	assetInfo := database.SearchAssetByOrderId(order.Db, orderUpdate.Id)
// 	fmt.Println(assetInfo)

// 	if orderUpdate.Id != c.Params("id") {
// 		return c.Status(400).JSON(&fiber.Map{
// 			"success": false,
// 			"message": "The order " + orderUpdate.Id +
// 				" from the body request is different from the " +
// 				c.Params("id"),
// 		})
// 	}

// 	orderUpdate.Country = assetInfo[0].AssetType.Country
// 	message := orderVerification(orderUpdate)
// 	if message != "" {
// 		return c.Status(400).JSON(&fiber.Map{
// 			"success": false,
// 			"message": message,
// 		})
// 	}

// 	orderUpdateReturn := database.UpdateOrderFromUser(order.Db,
// 		orderUpdate, userId.String())
// 	if orderUpdateReturn == nil {
// 		return c.Status(404).JSON(&fiber.Map{
// 			"success": false,
// 			"message": "The order Id of " + orderUpdate.Id +
// 				" does not exist for your user.",
// 		})
// 	}

// 	if err := c.JSON(&fiber.Map{
// 		"success": true,
// 		"order":   orderUpdateReturn,
// 		"message": "Order updated successfully",
// 	}); err != nil {
// 		return c.Status(500).JSON(&fiber.Map{
// 			"success": false,
// 			"message": err,
// 		})
// 	}

// 	return err
// }
