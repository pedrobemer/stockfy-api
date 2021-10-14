package fiberHandlers

import (
	"fmt"
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
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
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	httpStatusCode, orderCreated, err := order.LogicApi.ApiCreateOrder(
		orderInserted.Symbol, orderInserted.Country, orderInserted.OrderType,
		orderInserted.Quantity, orderInserted.Price, orderInserted.Currency,
		orderInserted.Brokerage, orderInserted.Date, userId.String())

	if httpStatusCode == 400 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if httpStatusCode == 404 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAssetSymbolUser.Error(),
			"code":    404,
		})
	}

	if httpStatusCode == 500 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
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
			"message": entity.ErrInvalidApiInternalError.Error(),
			"error":   err.Error(),
		})
	}

	return err

}

func (order *OrderApi) GetOrdersFromAssetUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	httpStatusCode, ordersInfo, err := order.LogicApi.ApiGetOrdersFromAssetUser(
		c.Query("symbol"), userId.String())

	if httpStatusCode == 400 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if httpStatusCode == 404 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAssetSymbolUser.Error(),
			"code":    404,
		})
	}

	if httpStatusCode == 500 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
		})
	}

	orderApiReturn := presenter.ConvertOrderToApiReturn(ordersInfo)

	if err := c.JSON(&fiber.Map{
		"success":    true,
		"ordersInfo": orderApiReturn,
		"message":    "Orders returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err
}

func (order *OrderApi) DeleteOrderFromUser(c *fiber.Ctx) error {
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	deletedOrderId, err := order.ApplicationLogic.OrderApp.DeleteOrdersFromUser(
		c.Params("id"), userId.String())
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if deletedOrderId == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiOrderId.Error(),
			"code":    404,
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"order":   deletedOrderId,
		"message": "Order deleted successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
		})
	}

	return err
}

func (order *OrderApi) UpdateOrderFromUser(c *fiber.Ctx) error {
	var err error

	var orderUpdate presenter.OrderBody
	if err := c.BodyParser(&orderUpdate); err != nil {
		fmt.Println(err)
	}

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	httpStatusCode, updatedOrder, err := order.LogicApi.ApiUpdateOrdersFromUser(
		c.Params("id"), userId.String(), orderUpdate.OrderType,
		orderUpdate.Price, orderUpdate.Quantity, orderUpdate.Date,
		orderUpdate.Brokerage)

	if httpStatusCode == 400 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if httpStatusCode == 404 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiOrderId.Error(),
			"code":    404,
		})
	}

	if httpStatusCode == 500 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
		})
	}

	orderApiReturn := presenter.ConvertSingleOrderToApiReturn(*updatedOrder)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"order":   orderApiReturn,
		"message": "Order updated successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiInternalError.Error(),
			"error":   err.Error(),
			"code":    500,
		})
	}

	return err
}
