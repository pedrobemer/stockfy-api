package handlers

import (
	"fmt"
	"stockfyApi/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func GetAsset(c *fiber.Ctx) error {

	var symbolQuery []database.AssetQueryReturn
	var err error

	symbolQuery = database.SearchAsset(database.DBpool, c.Params("symbol"), "")
	if symbolQuery == nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "SearchAsset: The symbol " + c.Params("symbol") +
				" is not registered in your database.",
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": symbolQuery,
		"message":   "Asset information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func GetAssetWithOrders(c *fiber.Ctx) error {

	var symbolQuery []database.AssetQueryReturn
	var err error
	var orderType string

	if c.Params("*") == "" {
		orderType = "ONLYORDERS"
	} else if c.Params("*") == "%26with-orders-info" {
		orderType = "ALL"
	} else if c.Params("*") == "%26only-orders-info" {
		orderType = "ONLYINFO"
	} else {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please see our documentation to properly execute requests for our API.",
		})
	}

	symbolQuery = database.SearchAsset(database.DBpool, c.Params("symbol"),
		orderType)
	if symbolQuery == nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "SearchAsset: The symbol " + c.Params("symbol") +
				" is not registered in your database.",
		})
	}

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": symbolQuery,
		"message":   "Asset information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func PostAsset(c *fiber.Ctx) error {

	var assetInsert database.AssetInsert
	var err error
	var assetTypeId string
	var sectorId string

	if err := c.BodyParser(&assetInsert); err != nil {
		fmt.Println(err)
	}
	fmt.Println(assetInsert)

	var condAssetExist = "symbol='" + assetInsert.Symbol + "'"
	assetExist := database.VerifyRowExistence(*database.DBpool, "asset",
		condAssetExist)
	fmt.Println(assetExist)
	if assetExist {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Asset already exist in our database",
		})
	}

	assetTypeQuery, _ := database.FetchAssetType(*database.DBpool, true,
		assetInsert.AssetType, assetInsert.Country)
	assetTypeId = assetTypeQuery[0].Id

	sectorQuery, _ := database.CreateSector(*database.DBpool, assetInsert.Sector)
	if sectorQuery != nil {
		sectorId = sectorQuery[0].Id
	}

	fmt.Println(assetTypeId, sectorId)

	symbolInsert := database.CreateAsset(*database.DBpool, assetInsert,
		assetTypeId, sectorId)

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": symbolInsert,
		"message":   "Asset creation was sucessful",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
