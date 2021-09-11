package handlers

import (
	"fmt"
	"reflect"
	"stockfyApi/alphaVantage"
	"stockfyApi/commonTypes"
	"stockfyApi/database"
	"stockfyApi/finnhub"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type AssetApi struct {
	Db database.PgxIface
}

func (asset *AssetApi) GetAsset(c *fiber.Ctx) error {

	var symbolQuery []database.AssetQueryReturn
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	symbolQuery, _ = database.SearchAssetByUser(asset.Db, c.Params("symbol"),
		userId.String(), "")
	if symbolQuery == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "SearchAssetByUser: The symbol " + c.Params("symbol") +
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

func (asset *AssetApi) GetAssetWithOrders(c *fiber.Ctx) error {

	var symbolQuery []database.AssetQueryReturn
	var err error
	var orderType string

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("withInfo") == "" && c.Query("onlyInfo") == "" {
		orderType = "ONLYORDERS"
	} else if c.Query("withInfo") == "true" {
		orderType = "ALL"
	} else if c.Query("onlyInfo") == "true" {
		orderType = "ONLYINFO"
	} else {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong REST API. Please see our documentation to properly execute requests for our API.",
		})
	}

	symbolQuery, _ = database.SearchAssetByUser(asset.Db, c.Params("symbol"),
		userId.String(), orderType)
	if symbolQuery == nil {
		return c.Status(404).JSON(&fiber.Map{
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

func (asset *AssetApi) GetAssetsFromAssetType(c *fiber.Ctx) error {

	var assetTypeQuery []database.AssetTypeApiReturn
	var err error
	var withOrdersInfo bool

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("type") == "" || c.Query("country") == "" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong Query. Please read our REST API description.",
		})
	}

	if c.Query("ordersInfo") == "true" {
		withOrdersInfo = true
	} else {
		withOrdersInfo = false
	}

	assetTypeQuery = database.SearchAssetsPerAssetType(asset.Db,
		c.Query("type"), c.Query("country"), userId.String(), withOrdersInfo)
	if assetTypeQuery == nil {
		message := "SearchAssetsPerAssetType: There is no asset registered as " +
			c.Query("type") + " from country " + c.Query("country")
		return c.Status(404).JSON(&fiber.Map{
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

func (asset *AssetApi) PostAsset(c *fiber.Ctx) error {

	var assetInsert database.AssetInsert
	var err error
	var apiType string

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	userInfoDb, err := database.SearchUser(asset.Db, userId.String())
	if userInfoDb[0].Type != "admin" {
		return c.Status(405).JSON(&fiber.Map{
			"success": false,
			"message": "User is not authorized to create Assets",
		})
	}

	if err := c.BodyParser(&assetInsert); err != nil {
		fmt.Println(err)
	}
	fmt.Println(assetInsert)

	var condAssetExist = "symbol='" + assetInsert.Symbol + "'"
	assetExist := database.VerifyRowExistence(asset.Db, "asset",
		condAssetExist)
	fmt.Println(assetExist)
	if assetExist {
		return c.Status(409).JSON(&fiber.Map{
			"success": false,
			"message": "Asset already exist in our database",
		})
	}

	if assetInsert.Country == "BR" {
		apiType = "Alpha"
	} else {
		apiType = "Finnhub"
	}
	symbolInsert, _ := asset.assetVerification(assetInsert.Symbol,
		assetInsert.Country, apiType)

	if symbolInsert.Symbol == "" {
		return c.Status(400).JSON(&fiber.Map{
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

func (asset *AssetApi) assetVerification(symbol string, country string, apiType string) (
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
	assetTypeQuery, err := database.FetchAssetType(asset.Db,
		"SPECIFIC", symbolLookup.Type, country)
	fmt.Println(assetTypeQuery)
	if err != nil {
		panic(err)
	}
	assetTypeId = assetTypeQuery[0].Id

	if sectorName != "" {
		sectorInfo, _ := database.CreateSector(asset.Db, sectorName)
		fmt.Println(sectorInfo)
		Ids.SectorId = sectorInfo[0].Id
	}
	fmt.Println(Ids.SectorId, assetTypeId)

	assetInsert.Fullname = symbolLookup.Fullname
	assetInsert.Symbol = symbol
	assetInsert.Country = country
	assetInsert.AssetType = assetTypeQuery[0].Type
	symbolInserted = database.CreateAsset(asset.Db, assetInsert,
		assetTypeId, Ids.SectorId)
	Ids.AssetId = symbolInserted.Id
	fmt.Println(symbolInserted)

	return symbolInserted, Ids
}

func (asset *AssetApi) DeleteAsset(c *fiber.Ctx) error {
	var err error
	var assetInfo []database.AssetQueryReturn

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	userInfoDb, err := database.SearchUser(asset.Db, userId.String())
	if userInfoDb[0].Type != "admin" && c.Query("myUser") == "" {
		return c.Status(405).JSON(&fiber.Map{
			"success": false,
			"message": "User is not authorized to delete Assets",
		})
	}

	if c.Query("myUser") == "" {

		assetInfo, err = database.SearchAsset(asset.Db, c.Params("symbol"))
		if assetInfo[0].Symbol == "" {
			return c.Status(404).JSON(&fiber.Map{
				"success": false,
				"message": "The Asset " + c.Query("symbol") + " does not exist in " +
					"the asset table. Please provide a valid symbol.",
			})
		}

		database.DeleteAssetUserRelationByAsset(asset.Db, assetInfo[0].Id)

		ordersId := database.DeleteOrdersFromAsset(asset.Db, assetInfo[0].Id)

		assetInfo = database.DeleteAsset(asset.Db, assetInfo[0].Id)

		assetInfo[0].OrdersList = ordersId

	} else if c.Query("myUser") == "true" {
		assetInfo, _ = database.SearchAsset(asset.Db, c.Params("symbol"))
		if assetInfo[0].Symbol == "" {
			return c.Status(404).JSON(&fiber.Map{
				"success": false,
				"message": "The Asset " + c.Query("symbol") + " does not exist in " +
					" the Asset table. Please provide a valid symbol.",
			})
		}

		database.DeleteOrdersFromAssetUser(asset.Db, assetInfo[0].Id,
			userId.String())

		assetUserReturn, _ := database.DeleteAssetUserRelation(asset.Db,
			assetInfo[0].Id, userId.String())
		if assetUserReturn[0].AssetId == "" {
			return c.Status(404).JSON(&fiber.Map{
				"success": false,
				"message": "The Asset " + c.Query("symbol") + " does not exist in " +
					"your Asset table. Please provide a valid symbol.",
			})
		}

	} else {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Unknown value for myUser variable in the REST API",
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"asset":   assetInfo,
		"message": "Asset was deleted successfuly",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
