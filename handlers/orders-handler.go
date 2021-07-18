package handlers

import (
	"fmt"
	"stockfyApi/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func PostOrder(c *fiber.Ctx) error {

	var orderInsert database.OrderBodyPost
	if err := c.BodyParser(&orderInsert); err != nil {
		fmt.Println(err)
	}
	fmt.Println(orderInsert)

	var assetId string
	var assetTypeId string
	var sectorId string
	var brokerageId string
	var orderReturn database.OrderApiReturn
	var err error
	var assetExist bool

	var condAssetExist = "symbol='" + orderInsert.Symbol + "'"
	assetExist = database.VerifyRowExistence(*database.DBpool, "asset",
		condAssetExist)

	if !assetExist {
		var assetInsert database.AssetInsert

		assetTypeQuery, err := database.FetchAssetType(*database.DBpool,
			"SPECIFIC", orderInsert.AssetType, orderInsert.Country)
		if err != nil {
			panic(err)
		}
		assetTypeId = assetTypeQuery[0].Id

		if orderInsert.Sector != "" {
			sectorInfo, _ := database.CreateSector(*database.DBpool,
				orderInsert.Sector)
			sectorId = sectorInfo[0].Id
		}

		assetInsert.Fullname = orderInsert.Fullname
		assetInsert.Symbol = orderInsert.Symbol
		assetInsert.Country = orderInsert.Country
		assetInsert.AssetType = assetTypeQuery[0].Type
		symbolInserted := database.CreateAsset(*database.DBpool, assetInsert,
			assetTypeId, sectorId)
		assetId = symbolInserted.Id
		fmt.Println(symbolInserted)
	} else {
		symbolQuery := database.SearchAsset(database.DBpool,
			orderInsert.Symbol, "")
		assetId = symbolQuery[0].Id
	}

	brokerageReturn, _ := database.FetchBrokerage(*database.DBpool, "SINGLE",
		orderInsert.Brokerage)
	brokerageId = brokerageReturn[0].Id

	fmt.Println(assetId, brokerageId)

	orderReturn = database.CreateOrder(*database.DBpool, orderInsert, assetId,
		brokerageId)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"orders":  orderReturn,
		"message": "Order registered successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}
