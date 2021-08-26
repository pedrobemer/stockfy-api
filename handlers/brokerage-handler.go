package handlers

import (
	"stockfyApi/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type BrokerageApi struct {
	Db database.PgxIface
}

func (brokerage *BrokerageApi) GetBrokerageFirms(c *fiber.Ctx) error {

	var brokerageQuery []database.BrokerageApiReturn
	var specificFetch string
	// var country string
	var err error

	if c.Query("country") == "" {
		specificFetch = "ALL"
	} else if c.Query("country") == "US" || c.Query("country") == "BR" {
		specificFetch = "COUNTRY"
	} else {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please see our documentation to understand how to use our API.",
		})
	}

	brokerageQuery, _ = database.FetchBrokerage(brokerage.Db, specificFetch,
		c.Query("country"))

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"brokerage": brokerageQuery,
		"message":   "Brokerage information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func (brokerage *BrokerageApi) GetBrokerageFirm(c *fiber.Ctx) error {

	var brokerageQuery []database.BrokerageApiReturn
	var err error

	brokerageQuery, _ = database.FetchBrokerage(brokerage.Db, "SINGLE",
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
		"brokerage": brokerageQuery,
		"message":   "Brokerage information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
