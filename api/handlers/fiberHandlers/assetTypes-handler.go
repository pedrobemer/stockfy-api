package fiberHandlers

import (
	"fmt"
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
	LogicApi         logicApi.Application
}

func (assetType *AssetTypeApi) GetAssetTypes(c *fiber.Ctx) error {

	ordersResume := false

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	if c.Query("ordersResume") == "true" {
		ordersResume = true
	} else if c.Query("ordersResume") != "" && c.Query("ordersResume") != "false" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   entity.ErrInvalidApiQueryWithOrderResume.Error(),
			"code":    400,
		})
	}

	httpStatusCode, searchedAssetType, err := assetType.LogicApi.ApiAssetsPerAssetType(
		c.Query("type"), c.Query("country"), ordersResume, userId.String())
	if err != nil {
		return c.Status(httpStatusCode).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiRequest.Error(),
			"error":   err.Error(),
			"code":    400,
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
		"message":   "Asset type returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
