package handlers

import (
	"fmt"
	"stockfyApi/alphaVantage"
	"stockfyApi/commonTypes"
	"stockfyApi/database"
	"stockfyApi/finnhub"

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
		"success": true,
		"asset":   symbolQuery,
		"message": "Asset information returned successfully",
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

	if c.Query("withInfo") == "" && c.Query("onlyInfo") == "" {
		orderType = "ONLYORDERS"
	} else if c.Query("withInfo") == "true" {
		orderType = "ALL"
	} else if c.Query("onlyInfo") == "true" {
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
		"success": true,
		"asset":   symbolQuery,
		"message": "Asset information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func GetAssetsFromAssetType(c *fiber.Ctx) error {

	var assetTypeQuery []database.AssetTypeApiReturn
	var err error
	var withOrdersInfo bool

	if c.Query("type") == "" || c.Query("country") == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong Query. Please read our REST API description.",
		})
	}

	if c.Query("ordersInfo") == "true" {
		withOrdersInfo = true
	} else {
		withOrdersInfo = false
	}

	assetTypeQuery = database.SearchAssetsPerAssetType(*database.DBpool,
		c.Query("type"), c.Query("country"), withOrdersInfo)
	if assetTypeQuery == nil {
		message := "SearchAssetsPerAssetType: There is no asset registered as " +
			c.Query("type") + " from country " + c.Query("country")
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

func PostAsset(c *fiber.Ctx) error {

	var assetInsert database.AssetInsert
	var err error
	var apiType string

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

	if assetInsert.Country == "BR" {
		apiType = "Alpha"
	} else {
		apiType = "Finnhub"
	}
	symbolInsert, _ := assetVerification(assetInsert.Symbol, assetInsert.Country,
		apiType)

	if symbolInsert.Symbol == "" {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "The Asset " + assetInsert.Symbol + " from country " +
				assetInsert.Country + " does not exist",
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"asset":   symbolInsert,
		"message": "Asset creation was sucessful",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func assetVerification(symbol string, country string, apiType string) (
	database.AssetApiReturn, DatabaseId) {
	var assetTypeId string
	var Ids DatabaseId
	var assetInsert database.AssetInsert
	var searchSymbol string
	var sectorName string
	var symbolInserted database.AssetApiReturn
	var symbolLookup commonTypes.SymbolLookup
	var companyOverview alphaVantage.CompanyOverview
	var companyProfile finnhub.CompanyProfile2

	if country == "BR" {
		searchSymbol = symbol + ".SA"
	} else {
		searchSymbol = symbol
	}

	// ALPHA VANTAGE COMPLETE VERIFICATION
	if apiType == "Alpha" {
		symbolLookupAlpha := alphaVantage.VerifySymbolAlpha(searchSymbol)
		symbolLookup = alphaVantage.ConvertSymbolLookup(symbolLookupAlpha)
		if symbolLookup.Symbol == "" {
			return symbolInserted, Ids
		}
		fmt.Println(symbolLookup)

		if country == "BR" && symbolLookup.Type == "ETF" {
			symbolLookup.Type = "FII"
			for _, validEtf := range alphaVantage.ListValidBrETF {
				if symbolLookup.Symbol == validEtf {
					symbolLookup.Type = "ETF"
				}
			}
		} else if country != "BR" && symbolLookup.Type == "Equity" {
			companyOverview = alphaVantage.CompanyOverviewAlpha(searchSymbol)
			if companyOverview["Industry"] == "REAL ESTATE INVESTMENT TRUSTS" {
				symbolLookup.Type = "REIT"
			} else {
				symbolLookup.Type = "STOCK"
			}
		} else if symbolLookup.Type == "Equity" {
			symbolLookup.Type = "STOCK"
		}
		fmt.Println(symbolLookup)
	} else {
		// FINNHUB COMPLETE VERIFICATION
		symbolLookupFinnhub := finnhub.VerifySymbolFinnhub(searchSymbol)
		symbolLookup = finnhub.ConvertSymbolLookup(symbolLookupFinnhub)
		if symbolLookup.Symbol == "" {
			return symbolInserted, Ids
		}
	}

	fmt.Println(symbolLookup)
	if symbolLookup.Type == "STOCK" {
		companyProfile = finnhub.CompanyProfile2Finnhub(searchSymbol)
		sectorName = companyProfile.FinnhubIndustry
		fmt.Println(companyProfile)
	}

	fmt.Println(symbolLookup.Type, country)
	assetTypeQuery, err := database.FetchAssetType(*database.DBpool,
		"SPECIFIC", symbolLookup.Type, country)
	fmt.Println(assetTypeQuery)
	if err != nil {
		panic(err)
	}
	assetTypeId = assetTypeQuery[0].Id

	if sectorName != "" {
		sectorInfo, _ := database.CreateSector(*database.DBpool, sectorName)
		fmt.Println(sectorInfo)
		Ids.SectorId = sectorInfo[0].Id
	}
	fmt.Println(Ids.SectorId, assetTypeId)

	assetInsert.Fullname = symbolLookup.Fullname
	assetInsert.Symbol = symbol
	assetInsert.Country = country
	assetInsert.AssetType = assetTypeQuery[0].Type
	symbolInserted = database.CreateAsset(*database.DBpool, assetInsert,
		assetTypeId, Ids.SectorId)
	Ids.AssetId = symbolInserted.Id
	fmt.Println(symbolInserted)

	return symbolInserted, Ids
}
