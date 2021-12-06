package fiberHandlers

import (
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type AssetTypeApi struct {
	ApplicationLogic usecases.Applications
	LogicApi         logicApi.UseCases
}

func (assetType *AssetTypeApi) GetAssetTypes(c *fiber.Ctx) error {

	ordersResume := false
	withPrice := false

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("ordersResume") == "true" {
		ordersResume = true
	} else if c.Query("ordersResume") != "" && c.Query("ordersResume") != "false" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiQueryWithOrderResume.Error(),
			"code":    400,
		})
	}

	if c.Query("withPrice") == "true" {
		withPrice = true
	} else if c.Query("withPrice") != "" && c.Query("withPrice") != "false" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   entity.ErrInvalidApiQueryWithPrice.Error(),
			"code":    400,
		})
	}

	httpStatusCode, searchedAssetType, err := assetType.LogicApi.ApiAssetsPerAssetType(
		c.Query("type"), c.Query("country"), ordersResume, withPrice,
		userId.String())
	if err != nil {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrMessageApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
		})
	}

	assetApiReturn := presenter.ConvertAssetTypeToApiReturn(searchedAssetType.Id,
		searchedAssetType.Type, searchedAssetType.Name, searchedAssetType.Country)
	sliceAssetsApiReturn := presenter.ConvertArrayAssetApiReturn(
		searchedAssetType.Assets)
	assetApiReturn.Assets = sliceAssetsApiReturn

	err = c.JSON(&fiber.Map{
		"success":   true,
		"assetType": assetApiReturn,
		"message":   "Asset type returned successfully",
	})

	return err
}
