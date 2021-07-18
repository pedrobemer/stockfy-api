package handlers

import (
	"stockfyApi/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func GetBrokerageFirms(c *fiber.Ctx) error {

	var brokerageQuery []database.BrokerageApiReturn
	var specificFetch string
	var country string
	var err error

	// if c.Params("searchType") == "SINGLE" && c.Params("name") == "" {
	// 	return c.SendString("The " + c.Params("searchType") + " needs the name field to be populated with a name of some brokerage firm. For example: /brokerage/SINGLE/<BROKERAGE FIRM NAME>.")
	// } else if c.Params("searchType") == "COUNTRY" &&
	// 	(c.Params("name") != "BR" && c.Params("name") != "US") {
	// 	return c.SendString("The " + c.Params("searchType") + " needs the name field to be populated with the country code BR or US. For example: /brokerage/COUNTRY/BR.")
	// } else if c.Params("searchType") != "SINGLE" &&
	// 	c.Params("searchType") != "ALL" &&
	// 	c.Params("searchType") != "COUNTRY" {
	// 	return c.SendString("Wrong Path. Please see the documentation from our REST API.")
	// }
	if c.Params("*") == "" {
		specificFetch = "ALL"
	} else if c.Params("*") == "%26country=US" {
		specificFetch = "COUNTRY"
		country = "US"
	} else if c.Params("*") == "%26country=BR" {
		specificFetch = "COUNTRY"
		country = "BR"
	} else {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please see our documentation to understand how to use our API.",
		})
	}

	brokerageQuery, _ = database.FetchBrokerage(*database.DBpool, specificFetch,
		country)

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": brokerageQuery,
		"message":   "Brokerage information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func GetBrokerageFirm(c *fiber.Ctx) error {

	var brokerageQuery []database.BrokerageApiReturn
	var err error

	brokerageQuery, _ = database.FetchBrokerage(*database.DBpool, "SINGLE",
		c.Params("name"))
	if brokerageQuery == nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Fetch Brokerage: The brokerage firm " +
				c.Params("name") + " does not exist in our database.",
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": brokerageQuery,
		"message":   "Brokerage information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
