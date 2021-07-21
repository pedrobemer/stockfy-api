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

	var symbolLookupInfo = finnhub.VerifySymbolFinnhub(c.Query("symbol"))
	symbolLookupUnique = finnhub.ConvertSymbolLookup(symbolLookupInfo)

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

func GetCompanyProfile2Finnhub(c *fiber.Ctx) error {
	var err error
	var message string

	if c.Query("symbol") == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	companyProfile2 := finnhub.CompanyProfile2Finnhub(c.Query("symbol"))

	if companyProfile2.Ticker == "" {
		message = "Company Profile 2 via Finnhub returned successfully, but " +
			"there is no company with symbol " + c.Query("symbol")
	} else {
		message = "Company Profile 2 via Finnhub returned successfully"
	}

	if err := c.JSON(&fiber.Map{
		"success":         true,
		"companyProfile2": companyProfile2,
		"message":         message,
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
