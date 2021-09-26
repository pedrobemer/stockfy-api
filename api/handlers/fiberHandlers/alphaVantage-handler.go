package handlers

import (
	"stockfyApi/alphaVantage"

	"github.com/gofiber/fiber/v2"
)

type AlphaVantageApi struct{}

type httpErrorResp struct {
	Message    string
	Success    bool
	StatusCode int
}

func (alpha *AlphaVantageApi) GetSymbolAlphaVantage(c *fiber.Ctx) error {
	var err error
	var message string

	httpResp := apiVerification(c)
	if httpResp.Message != "" {
		return c.Status(httpResp.StatusCode).JSON(&fiber.Map{
			"success": httpResp.Success,
			"message": httpResp.Message,
		})
	}

	searchSymbol, country := convertSymbol(c.Query("symbol"), c.Query("country"))

	symbolLookup := alphaVantage.VerifySymbolAlpha(searchSymbol)
	symbolLookupUnique := alphaVantage.ConvertSymbolLookup(symbolLookup)

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
	var message string

	httpResp := apiVerification(c)
	if httpResp.Message != "" {
		return c.Status(httpResp.StatusCode).JSON(&fiber.Map{
			"success": httpResp.Success,
			"message": httpResp.Message,
		})
	}

	searchSymbol, country := convertSymbol(c.Query("symbol"), c.Query("country"))

	companyOverview := alphaVantage.CompanyOverviewAlpha(searchSymbol)

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
	var message string

	httpResp := apiVerification(c)
	if httpResp.Message != "" {
		return c.Status(httpResp.StatusCode).JSON(&fiber.Map{
			"success": httpResp.Success,
			"message": httpResp.Message,
		})
	}

	searchSymbol, country := convertSymbol(c.Query("symbol"), c.Query("country"))

	symbolPrice := alphaVantage.GetPriceAlphaVantage(searchSymbol)

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

func apiVerification(c *fiber.Ctx) httpErrorResp {

	if c.Query("symbol") == "" {
		return httpErrorResp{
			Message:    "Symbol not specified. Please read our REST API documentation",
			Success:    false,
			StatusCode: 400,
		}
	}

	if c.Query("country") != "BR" && c.Query("country") != "US" {
		return httpErrorResp{
			Message:    "Country not specified or it is invalid. Please read our REST API documentation.",
			Success:    false,
			StatusCode: 400,
		}
	}

	return httpErrorResp{}
}

func convertSymbol(symbol string, country string) (string, string) {
	var searchSymbol string

	if country == "BR" {
		searchSymbol = symbol + ".SA"
	} else {
		searchSymbol = symbol
	}

	return searchSymbol, country
}
