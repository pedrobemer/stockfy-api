package fiberHandlers

import (
	"fmt"
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
	LogicApi           logicApi.Application
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

func (asset *AssetApi) GetAssetsFromAssetType(c *fiber.Ctx) error {
	var withOrdersInfo bool

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("ordersInfo") == "true" {
		withOrdersInfo = true
	} else {
		withOrdersInfo = false
	}

	httpStatusCode, searchedAssetType, err := asset.LogicApi.ApiAssetsPerAssetType(
		c.Query("type"), c.Query("country"), withOrdersInfo, userId.String())
	if err != nil {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	assetApiReturn := presenter.ConvertAssetTypeToApiReturn(searchedAssetType.Id,
		searchedAssetType.Type, searchedAssetType.Name, searchedAssetType.Country)
	fmt.Println(assetApiReturn)
	sliceAssetsApiReturn := presenter.ConvertArrayAssetApiReturn(
		searchedAssetType.Assets)
	assetApiReturn.Assets = &sliceAssetsApiReturn

	if err := c.JSON(&fiber.Map{
		"success":   true,
		"assetType": assetApiReturn,
		"message":   "The asset type returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

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
	statusCode, assetCreated, err := asset.LogicApi.ApiAssetVerification(
		assetInsert.Symbol, assetInsert.Country)
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

func (asset *AssetApi) DeleteAsset(c *fiber.Ctx) error {

	myUser := false

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("myUser") == "true" {
		myUser = true
	}

	httpStatusCode, deletedAsset, err := asset.LogicApi.ApiDeleteAssets(myUser,
		userId.String(), c.Params("symbol"))
	if err != nil {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	deletedAssetApiReturn := presenter.ConvertAssetToApiReturn(deletedAsset.Id,
		*deletedAsset.Preference, deletedAsset.Fullname, deletedAsset.Symbol,
		deletedAsset.Sector.Name, deletedAsset.Sector.Id, deletedAsset.AssetType.Id,
		deletedAsset.AssetType.Type, deletedAsset.AssetType.Country,
		deletedAsset.AssetType.Name, nil, nil)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"asset":   deletedAssetApiReturn,
		"message": "Asset was deleted successfuly",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err
}
