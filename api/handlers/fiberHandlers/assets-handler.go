package fiberHandlers

import (
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type AssetApi struct {
	ApplicationLogic   usecases.Applications
	ExternalInterfaces externalapi.ThirdPartyInterfaces
}

func (asset *AssetApi) GetAsset(c *fiber.Ctx) error {

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	searchedAsset, err := asset.ApplicationLogic.AssetApp.SearchAssetByUser(
		c.Params("symbol"), userId.String(), false, false, true)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	assetApiReturn := presenter.ConvertAssetToApiReturn(searchedAsset.Id,
		*searchedAsset.Preference, searchedAsset.Fullname, searchedAsset.Symbol,
		searchedAsset.Sector.Name, searchedAsset.Sector.Id, searchedAsset.AssetType.Id,
		searchedAsset.AssetType.Type, searchedAsset.AssetType.Country,
		searchedAsset.AssetType.Name, nil, nil)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"asset":   assetApiReturn,
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

	withInfo := false
	onlyInfo := false

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("withInfo") == "true" {
		withInfo = true
	} else if c.Query("withInfo") == "" {
		withInfo = false
	} else {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
		})
	}

	if c.Query("onlyInfo") == "true" {
		onlyInfo = true
	} else if c.Query("onlyInfo") == "" {
		onlyInfo = false
	} else {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
		})
	}

	searchedAsset, err := asset.ApplicationLogic.AssetApp.SearchAssetByUser(
		c.Params("symbol"), userId.String(), withInfo, onlyInfo,
		false)
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	assetApiReturn := presenter.ConvertAssetToApiReturn(searchedAsset.Id,
		*searchedAsset.Preference, searchedAsset.Fullname, searchedAsset.Symbol,
		searchedAsset.Sector.Name, searchedAsset.Sector.Id, searchedAsset.AssetType.Id,
		searchedAsset.AssetType.Type, searchedAsset.AssetType.Country,
		searchedAsset.AssetType.Name, searchedAsset.OrdersList,
		searchedAsset.OrderInfo)

	presenter.ConvertOrderToApiReturn(searchedAsset.OrdersList)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"asset":   assetApiReturn,
		"message": "Asset information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

// func (asset *AssetApi) GetAssetsFromAssetType(c *fiber.Ctx) error {

// 	var assetTypeQuery []database.AssetTypeApiReturn
// 	var err error
// 	var withOrdersInfo bool

// 	userInfo := c.Context().Value("user")
// 	userId := reflect.ValueOf(userInfo).FieldByName("userID")

// 	if c.Query("type") == "" || c.Query("country") == "" {
// 		return c.Status(400).JSON(&fiber.Map{
// 			"success": false,
// 			"message": "Wrong Query. Please read our REST API description.",
// 		})
// 	}

// 	if c.Query("ordersInfo") == "true" {
// 		withOrdersInfo = true
// 	} else {
// 		withOrdersInfo = false
// 	}

// 	assetTypeQuery = database.SearchAssetsPerAssetType(asset.Db,
// 		c.Query("type"), c.Query("country"), userId.String(), withOrdersInfo)
// 	if assetTypeQuery == nil {
// 		message := "SearchAssetsPerAssetType: There is no asset registered as " +
// 			c.Query("type") + " from country " + c.Query("country")
// 		return c.Status(404).JSON(&fiber.Map{
// 			"success": false,
// 			"message": message,
// 		})
// 	}

// 	if err := c.JSON(&fiber.Map{
// 		"success":   true,
// 		"assetType": assetTypeQuery,
// 		"message":   "The asset type returned successfully",
// 	}); err != nil {
// 		return c.Status(500).JSON(&fiber.Map{
// 			"success": false,
// 			"message": err,
// 		})
// 	}

// 	return err

// }

func (asset *AssetApi) CreateAsset(c *fiber.Ctx) error {

	var assetInsert presenter.AssetBody
	var err error

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	// Verify if this is a Admin user. If not, this user is not authorized to
	// create an asset.
	searchedUser, _ := asset.ApplicationLogic.UserApp.SearchUser(userId.String())
	if searchedUser.Type != "admin" {
		return c.Status(405).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAuthorization.Error(),
		})
	}

	if err := c.BodyParser(&assetInsert); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong JSON in the Body",
		})
	}

	// Verify if this Asset is already in our database
	condAssetExist := "symbol='" + assetInsert.Symbol + "'"
	assetExist := asset.ApplicationLogic.DbVerificationApp.RowValidation("asset",
		condAssetExist)
	if assetExist {
		return c.Status(409).JSON(&fiber.Map{
			"success": false,
			"message": "Asset already exist in our database",
		})
	}

	// Verify if this asset exist in the US or BR stock market
	statusCode, assetCreated, err := logicApi.ApiAssetVerification(
		asset.ApplicationLogic, asset.ExternalInterfaces, assetInsert.Symbol,
		assetInsert.Country)
	if err != nil {
		return c.Status(statusCode).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"asset": presenter.ConvertAssetToApiReturn(assetCreated.Id,
			*assetCreated.Preference, assetCreated.Fullname,
			assetCreated.Symbol, "", "", "", "", "", "", nil, nil),
		"message": "Asset creation was sucessful",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err

}

// func (asset *AssetApi) DeleteAsset(c *fiber.Ctx) error {
// 	// var err error
// 	// var assetInfo []database.AssetQueryReturn

// 	userInfo := c.Context().Value("user")
// 	userId := reflect.ValueOf(userInfo).FieldByName("userID")

// 	if c.Query("myUser") == "" {

// 		searchedUser, _ := asset.ApplicationLogic.UserApp.SearchUser(
// 			userId.String())
// 		if searchedUser.Type != "admin" {
// 			return c.Status(405).JSON(&fiber.Map{
// 				"success": false,
// 				"message": entity.ErrInvalidApiAuthorization.Error(),
// 			})
// 		}

// 		// assetInfo, err = database.SearchAsset(asset.Db, c.Params("symbol"))
// 		assetInfo, err := asset.ApplicationLogic.AssetApp.SearchAsset(c.Params("symbol"))

// 		if err != nil {
// 			return c.Status(400).JSON(&fiber.Map{
// 				"success": false,
// 				"message": err.Error(),
// 			})
// 		}

// 		database.DeleteAssetUserRelationByAsset(asset.Db, assetInfo[0].Id)

// 		ordersId := database.DeleteOrdersFromAsset(asset.Db, assetInfo[0].Id)

// 		assetInfo = database.DeleteAsset(asset.Db, assetInfo[0].Id)

// 		assetInfo[0].OrdersList = ordersId

// 	} else if c.Query("myUser") == "true" {
// 		asset.ApplicationLogic.AssetApp.SearchAssetByUser(c.Params("symbol"),
// 			userId.String(), false, false, true)
// 		// assetInfo, _ = database.SearchAsset(asset.Db, c.Params("symbol"))
// 		// if assetInfo[0].Symbol == "" {
// 		// 	return c.Status(404).JSON(&fiber.Map{
// 		// 		"success": false,
// 		// 		"message": "The Asset " + c.Query("symbol") + " does not exist in " +
// 		// 			" the Asset table. Please provide a valid symbol.",
// 		// 	})
// 		// }

// 		database.DeleteOrdersFromAssetUser(asset.Db, assetInfo[0].Id,
// 			userId.String())

// 		assetUserReturn, _ := database.DeleteAssetUserRelation(asset.Db,
// 			assetInfo[0].Id, userId.String())
// 		if assetUserReturn[0].AssetId == "" {
// 			return c.Status(404).JSON(&fiber.Map{
// 				"success": false,
// 				"message": "The Asset " + c.Query("symbol") + " does not exist in " +
// 					"your Asset table. Please provide a valid symbol.",
// 			})
// 		}

// 	} else {
// 		return c.Status(400).JSON(&fiber.Map{
// 			"success": false,
// 			"message": "Unknown value for myUser variable in the REST API",
// 		})
// 	}

// 	if err := c.JSON(&fiber.Map{
// 		"success": true,
// 		"asset":   assetInfo,
// 		"message": "Asset was deleted successfuly",
// 	}); err != nil {
// 		return c.Status(500).JSON(&fiber.Map{
// 			"success": false,
// 			"message": err,
// 		})
// 	}

// 	return err
// }
