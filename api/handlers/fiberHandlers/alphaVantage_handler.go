package fiberHandlers

import (
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/asset"

	"github.com/gofiber/fiber/v2"
)

type AlphaVantageApi struct {
	ApplicationLogic usecases.Applications
	Api              asset.ExternalApiRepository
}

func (alpha *AlphaVantageApi) GetSymbol(c *fiber.Ctx) error {

	symbolLookup, err := alpha.ApplicationLogic.AssetApp.
		AssetVerificationExistence(c.Query("symbol"), c.Query("country"),
			alpha.Api)

	if err == entity.ErrInvalidAssetSymbol {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"code":    404,
		})
	}

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	err = c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": symbolLookup,
		"message":      "Symbol Lookup via Alpha Vantage returned successfully",
	})

	return err

}

func (alpha *AlphaVantageApi) GetSymbolPrice(c *fiber.Ctx) error {

	symbolPrice, err := alpha.ApplicationLogic.AssetApp.AssetVerificationPrice(
		c.Query("symbol"), c.Query("country"), alpha.Api)

	if err == entity.ErrInvalidAssetSymbol {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
			"code":    404,
		})
	}

	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	err = c.JSON(&fiber.Map{
		"success":     true,
		"symbolPrice": symbolPrice,
		"message":     "Symbol Price via Alpha Vantage returned successfully",
	})

	return err

}
