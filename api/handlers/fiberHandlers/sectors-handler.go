package fiberHandlers

import (
	"reflect"
	"stockfyApi/api/presenter"
	"stockfyApi/entity"
	"stockfyApi/usecases"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type SectorApi struct {
	ApplicationLogic usecases.Applications
}

func (sector *SectorApi) GetSector(c *fiber.Ctx) error {

	sectorInfo, err := sector.ApplicationLogic.SectorApp.SearchSectorByName(
		c.Params("sector"))
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if sectorInfo == nil {
		return c.Status(404).JSON(&fiber.Map{
			"success": false,
			"error":   entity.ErrInvalidSectorSearchName.Error(),
		})
	}

	sectorApiReturn := presenter.ConvertSectorToApiReturn(sectorInfo.Id, sectorInfo.Name)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"sector":  sectorApiReturn,
		"message": "Sector information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return err

}

func (sector *SectorApi) CreateSector(c *fiber.Ctx) error {
	var sectorBody presenter.SectorBody

	userInfo := c.Context().Value("user")
	userId := reflect.ValueOf(userInfo).FieldByName("userID")

	// Verify if this is a Admin user. If not, this user is not authorized to
	// create an asset.
	searchedUser, _ := sector.ApplicationLogic.UserApp.SearchUser(userId.String())
	if searchedUser.Type != "admin" {
		return c.Status(405).JSON(&fiber.Map{
			"success": false,
			"message": entity.ErrInvalidApiAuthorization.Error(),
		})
	}

	if err := c.BodyParser(&sectorBody); err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": "Wrong JSON in the Body",
		})
	}

	sectorCreated, err := sector.ApplicationLogic.SectorApp.CreateSector(
		sectorBody.Sector)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	sectorApiReturn := presenter.ConvertSectorToApiReturn(sectorCreated[0].Id,
		sectorCreated[0].Name)

	if err := c.JSON(&fiber.Map{
		"success": true,
		"sector":  sectorApiReturn,
		"message": "Created sector successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
