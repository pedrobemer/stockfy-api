package fiberHandlers

import (
	"stockfyApi/usecases"
	"stockfyApi/usecases/asset"

	"github.com/gofiber/fiber/v2"
)

type AlphaVantageApi struct {
	ApplicationLogic usecases.Applications
	Api              asset.ExternalApiRepository
}

func (alpha *AlphaVantageApi) GetSymbol(c *fiber.Ctx) error {
	var err error
	var message string

	symbolLookup, err := alpha.ApplicationLogic.AssetApp.
		AssetVerificationExistence(c.Query("symbol"), c.Query("country"),
			alpha.Api)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": symbolLookup,
		"message":      message,
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func (alpha *AlphaVantageApi) GetSymbolPrice(c *fiber.Ctx) error {
	var err error
	var message string

	symbolPrice, err := alpha.ApplicationLogic.AssetApp.AssetVerificationPrice(
		c.Query("symbol"), c.Query("country"), alpha.Api)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":     true,
		"symbolPrice": symbolPrice,
		"message":     message,
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
