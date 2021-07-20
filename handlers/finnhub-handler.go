package handlers

import (
	"stockfyApi/commonTypes"
	"stockfyApi/finnhub"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func GetSymbolFinnhub(c *fiber.Ctx) error {
	var err error
	var symbolLookupUnique commonTypes.SymbolLookup

	var symbolLookup = finnhub.VerifySymbolFinnhub(c.Query("symbol"))

	for _, s := range symbolLookup.Result {
		if s.Symbol == c.Query("symbol") {
			symbolLookupUnique = finnhub.ConvertSymbolLookup(s)
		}
	}

	if err := c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": symbolLookupUnique,
		"message":      "Symbol Lookup via Finnhub returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func GetSymbolPriceFinnhub(c *fiber.Ctx) error {
	var err error

	if c.Query("symbol") == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	symbolPrice := finnhub.GetPriceFinnhub(c.Query("symbol"))

	if err := c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": symbolPrice,
		"message":      "Symbol Lookup via Finnhub returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
