package handlers

import (
	"stockfyApi/alphaVantage"
	"stockfyApi/commonTypes"

	"github.com/gofiber/fiber/v2"
)

type AlphaVantageApi struct{}

func (alpha *AlphaVantageApi) GetSymbolAlphaVantage(c *fiber.Ctx) error {
	var err error

	var symbolLookupUnique commonTypes.SymbolLookup
	var message string

	if c.Query("symbol") == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	if c.Query("country") != "BR" && c.Query("country") != "US" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	var searchSymbol string
	var country string
	if c.Query("country") == "BR" {
		searchSymbol = c.Query("symbol") + ".SA"
		country = c.Query("country")
	} else {
		searchSymbol = c.Query("symbol")
		country = "US"
	}

	var symbolLookup = alphaVantage.VerifySymbolAlpha(searchSymbol)
	symbolLookupUnique = alphaVantage.ConvertSymbolLookup(symbolLookup)

	if symbolLookupUnique.Symbol == "" {
		message = "Symbol Lookup via Alpha Vantage did not find the symbol " +
			c.Query("symbol") + " in country " + country
	} else {
		message = "Symbol Lookup via Alpha Vantage returned successfully"
	}

	if err := c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": symbolLookupUnique,
		"message":      message,
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func (alpha *AlphaVantageApi) GetCompanyOverviewAlphaVantage(c *fiber.Ctx) error {
	var err error

	var companyOverview alphaVantage.CompanyOverview
	var message string

	if c.Query("symbol") == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	if c.Query("country") != "BR" && c.Query("country") != "US" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	var searchSymbol string
	var country string
	if c.Query("country") == "BR" {
		searchSymbol = c.Query("symbol") + ".SA"
		country = c.Query("country")
	} else {
		searchSymbol = c.Query("symbol")
		country = "US"
	}

	companyOverview = alphaVantage.CompanyOverviewAlpha(searchSymbol)

	if companyOverview["Symbol"] == "" {
		message = "Symbol Lookup via Alpha Vantage did not find the symbol " +
			c.Query("symbol") + " in country " + country
	} else {
		message = "Symbol Lookup via Alpha Vantage returned successfully"
	}

	if err := c.JSON(&fiber.Map{
		"success":      true,
		"symbolLookup": companyOverview,
		"message":      message,
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func (alpha *AlphaVantageApi) GetSymbolPriceAlphaVantage(c *fiber.Ctx) error {
	var err error
	var symbolPrice commonTypes.SymbolPrice
	var message string
	var country string

	if c.Query("symbol") == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	if c.Query("country") != "BR" && c.Query("country") != "US" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please read our documentation.",
		})
	}

	var searchSymbol string
	if c.Query("country") == "BR" {
		searchSymbol = c.Query("symbol") + ".SA"
		country = c.Query("country")
	} else {
		searchSymbol = c.Query("symbol")
		country = "US"
	}

	symbolPrice = alphaVantage.GetPriceAlphaVantage(searchSymbol)

	if symbolPrice.Symbol == "" {
		message = "Symbol Lookup via Alpha Vantage did not find the symbol " +
			c.Query("symbol") + " in country " + country
	} else {
		message = "Symbol Lookup via Alpha Vantage returned successfully"
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
