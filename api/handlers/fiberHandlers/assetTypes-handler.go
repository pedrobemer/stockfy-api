package handlers

import (
	"stockfyApi/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type AssetTypeApi struct {
	Db database.PgxIface
}

func (assetType *AssetTypeApi) GetAssetTypes(c *fiber.Ctx) error {

	var assetTypeQuery []database.AssetTypeApiReturn
	var err error

	var fetchType string
	if c.Query("type") == "" && c.Query("country") == "" {
		fetchType = ""
	} else if c.Query("type") == "" && c.Query("country") != "" {
		fetchType = "ONLYCOUNTRY"
	} else if c.Query("type") != "" && c.Query("country") == "" {
		fetchType = "ONLYTYPE"
	} else if c.Query("type") != "" && c.Query("country") != "" {
		fetchType = "SPECIFIC"
	}

	assetTypeQuery, err = database.FetchAssetType(assetType.Db,
		fetchType, c.Query("type"), c.Query("country"))
	if err != nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"message": "FetchAssetType: There is no " + c.Query("type") +
				" from country " + c.Query("country") +
				" with this specifications.",
		})
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
