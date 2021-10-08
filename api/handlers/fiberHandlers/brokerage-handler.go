package fiberHandlers

import (
	"stockfyApi/api/presenter"
	"stockfyApi/usecases"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type BrokerageApi struct {
	ApplicationLogic usecases.Applications
}

func (brokerage *BrokerageApi) GetBrokerageFirms(c *fiber.Ctx) error {

	var searchType string

	if c.Query("country") == "" {
		searchType = "ALL"
	} else {
		searchType = "COUNTRY"
	}

	brokerageFirms, err := brokerage.ApplicationLogic.BrokerageApp.
		SearchBrokerage(searchType, "", c.Query("country"))
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	brokerageFirmsApiReturn := presenter.ConvertArrayBrokerageToApiReturn(
		brokerageFirms)

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"brokerage": brokerageFirmsApiReturn,
		"message":   "Brokerage information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err

}

func (brokerage *BrokerageApi) GetBrokerageFirm(c *fiber.Ctx) error {

	brokerageInfo, err := brokerage.ApplicationLogic.BrokerageApp.SearchBrokerage(
		"SINGLE", c.Params("name"), "")
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	brokerageApiReturn := presenter.ConvertBrokerageToApiReturn(brokerageInfo[0].Id,
		brokerageInfo[0].Name, brokerageInfo[0].Country)

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"brokerage": brokerageApiReturn,
		"message":   "Brokerage information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err

}
