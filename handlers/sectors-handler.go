package handlers

import (
	"fmt"
	"stockfyApi/database"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

type SectorApi struct {
	Db database.PgxIface
}

func (sector *SectorApi) GetAllSectors(c *fiber.Ctx) error {

	var sectorQuery []database.SectorApiReturn
	var err error
	sectorQuery, err = database.FetchSectorByName(sector.Db, "ALL")
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"sector":  sectorQuery,
		"message": "All sectors returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func (sector *SectorApi) GetSector(c *fiber.Ctx) error {

	var sectorQuery []database.SectorApiReturn
	var err error

	if c.Params("sector") == "ALL" {
		return c.Status(400).JSON(&fiber.Map{
			"success": false,
			"error":   "Unauthorized Sector Search",
		})
	}

	sectorQuery, err = database.FetchSectorByName(sector.Db, c.Params("sector"))
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"sector":  sectorQuery,
		"message": "Sector information returned successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err

}

func (sector *SectorApi) PostSector(c *fiber.Ctx) error {
	var sectorBodyPost database.SectorBodyPost
	var err error

	if err := c.BodyParser(&sectorBodyPost); err != nil {
		fmt.Println(err)
	}

	var sectorInsert []database.SectorApiReturn
	sectorInsert, err = database.CreateSector(sector.Db, sectorBodyPost.Sector)
	if err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	if err := c.JSON(&fiber.Map{
		"success": true,
		"sector":  sectorInsert,
		"message": "Created sector successfully",
	}); err != nil {
		return c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	return err
}
