package handlers

import (
	"fmt"
	"stockfyApi/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func GetAllAssetTypes(c *fiber.Ctx) error {

	var assetTypeQuery []database.AssetTypeApiReturn
	var err error

	var specificFetch = false

	assetTypeQuery, err = database.FetchAssetType(*database.DBpool, specificFetch)
	if err != nil {
		panic(err)
	}

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": assetTypeQuery,
		"message":   "All asset types returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func GetAssetTypesPerTypeAndCountry(c *fiber.Ctx) error {

	var assetTypeQuery []database.AssetTypeApiReturn
	var err error

	var specificFetch = true
	fmt.Println("TEST", c.Params("country"), c.Params("type"))
	assetTypeQuery, err = database.FetchAssetType(*database.DBpool, specificFetch,
		c.Params("type"), c.Params("country"))
	fmt.Println(err)
	if assetTypeQuery == nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": assetTypeQuery,
		"message":   "The asset type returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func GetAssetTypesWithAssets(c *fiber.Ctx) error {

	var assetTypeQuery []database.AssetTypeApiReturn
	var err error
	var withOrdersInfo bool

	fmt.Println(c.Params("*"))
	if c.Params("*") != "%26withOrdersInfo" && c.Params("*") != "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong Query. Please read our REST API description.",
		})
	} else if c.Params("*") == "%26withOrdersInfo" {
		withOrdersInfo = true
	} else {
		withOrdersInfo = false
	}
	assetTypeQuery = database.SearchAssetsPerAssetType(*database.DBpool,
		c.Params("type"), c.Params("country"), withOrdersInfo)
	if assetTypeQuery == nil {
		message := "SearchAssetsPerAssetType: There is no asset registered as " +
			c.Params("type") + " from country " + c.Params("country")
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": message,
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": assetTypeQuery,
		"message":   "The asset type returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
