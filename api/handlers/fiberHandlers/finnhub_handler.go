package fiberHandlers

import (
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/asset"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type FinnhubApi struct {
	ApplicationLogic usecases.Applications
	Api              asset.ExternalApiRepository
}

func (finn *FinnhubApi) GetSymbol(c *fiber.Ctx) error {
	var err error

	symbolLookup, err := finn.ApplicationLogic.AssetApp.
		AssetVerificationExistence(c.Query("symbol"), c.Query("country"),
			finn.Api)

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
		"message":      "Symbol Lookup via Finnhub returned successfully",
	})

	return err

}

func (finn *FinnhubApi) GetSymbolPrice(c *fiber.Ctx) error {
	var err error

	symbolPrice, err := finn.ApplicationLogic.AssetApp.AssetVerificationPrice(
		c.Query("symbol"), c.Query("country"), finn.Api)

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
		"message":     "Symbol Price via Finnhub returned successfully",
	})

	return err

}
