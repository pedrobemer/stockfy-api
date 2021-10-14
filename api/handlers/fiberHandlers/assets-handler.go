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
	LogicApi           logicApi.Application
}

func (asset *AssetApi) GetAsset(c *fiber.Ctx) error {

	withOrders := false
	withOrderResume := false

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("withOrders") == "true" {
		withOrders = true
	} else if c.Query("withOrders") == "" || c.Query("withOrders") == "false" {
		withOrders = false
	} else {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   entity.ErrInvalidApiQueryWithOrders.Error(),
			"code":    400,
		})
	}

	if c.Query("withOrderResume") == "true" {
		withOrderResume = true
	} else if c.Query("withOrderResume") == "" || c.Query("withOrderResume") == "false" {
		withOrderResume = false
	} else {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   entity.ErrInvalidApiQueryWithOrderResume.Error(),
			"code":    400,
		})
	}

	searchedAsset, err := asset.ApplicationLogic.AssetApp.SearchAssetByUser(
		c.Params("symbol"), userId.String(), withOrders, withOrderResume)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	if searchedAsset == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAssetSymbolUser.Error(),
			"code":    404,
		})
	}

	assetApiReturn := presenter.ConvertAssetToApiReturn(searchedAsset.Id,
		*searchedAsset.Preference, searchedAsset.Fullname, searchedAsset.Symbol,
		searchedAsset.Sector.Name, searchedAsset.Sector.Id, searchedAsset.AssetType.Id,
		searchedAsset.AssetType.Type, searchedAsset.AssetType.Country,
		searchedAsset.AssetType.Name, searchedAsset.OrdersList, searchedAsset.OrderInfo)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"asset":   assetApiReturn,
		"message": "Asset information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
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
		return c.Status(403).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAuthorization.Error(),
			"error":   entity.ErrInvalidApiUserAdminPrivilege.Error(),
			"code":    403,
		})
	}

	if err := c.BodyParser(&assetInsert); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   entity.ErrInvalidApiBody.Error(),
			"code":    400,
		})
	}

	// Verify if this Asset is already in our database
	condAssetExist := "symbol='" + assetInsert.Symbol + "'"
	assetExist := asset.ApplicationLogic.DbVerificationApp.RowValidation("asset",
		condAssetExist)
	if assetExist {
		return c.Status(403).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAuthorization.Error(),
			"error":   entity.ErrInvalidApiAssetSymbolExist.Error(),
			"code":    403,
		})
	}

	// Verify if this asset exist in the US or BR stock market
	statusCode, assetCreated, err := asset.LogicApi.ApiAssetVerification(
		assetInsert.Symbol, assetInsert.Country)

	if statusCode == 404 {
		return c.Status(statusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAssetSymbolUser.Error(),
			"code":    statusCode,
		})
	} else if statusCode == 400 {
		return c.Status(statusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   err.Error(),
			"code":    statusCode,
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

	if httpStatusCode == 403 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAuthorization.Error(),
			"error":   err.Error(),
			"code":    httpStatusCode,
		})
	}

	if httpStatusCode == 400 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   err.Error(),
			"code":    httpStatusCode,
		})
	}

	if httpStatusCode == 404 {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAssetSymbolUser.Error(),
			"code":    404,
		})
	}

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
