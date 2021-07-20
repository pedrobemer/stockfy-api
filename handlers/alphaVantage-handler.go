package handlers

import (
	"fmt"
	"stockfyApi/alphaVantage"
	"stockfyApi/commonTypes"

	"github.com/gofiber/fiber/v2"
)

func GetSymbolAlphaVantage(c *fiber.Ctx) error {
	var err error

	var symbolLookupUnique commonTypes.SymbolLookup
	// var symbolTypes = map[string]string{
	// 	"Equity": "STOCK",
	// 	"ETF":    "ETF",
	// 	"REIT":   "REIT",
	// 	"FII":    "FII",
	// }
	if c.Query("symbol") == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	var symbolLookup = alphaVantage.VerifySymbolAlpha(c.Query("symbol"))
	// fmt.Println(symbolLookup.BestMatches[0])
	// fmt.Println(symbolLookup.BestMatches[1])
	for index, s := range symbolLookup.BestMatches {
		if s["9. matchScore"] == "1.0000" {
			fmt.Println(index, "=>", s)
			symbolLookupUnique = alphaVantage.ConvertSymbolLookup(s)
		}
	}

	if err := c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": symbolLookupUnique,
		"message":      "Symbol Lookup via Alpha Vantage returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func GetSymbolPriceAlphaVantage(c *fiber.Ctx) error {
	var err error
	var symbolPrice commonTypes.SymbolPrice

	if c.Query("symbol") == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	symbolPrice = alphaVantage.GetPriceAlphaVantage(c.Query("symbol"))

	if err := c.JSON(&fiber.Map{
		"success":     true,
		"symbolPrice": symbolPrice,
		"message":     "Symbol Price via Alpha Vantage returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
