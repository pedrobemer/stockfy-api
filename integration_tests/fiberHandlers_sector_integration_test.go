package integration_tests

import (
	"encoding/json"
	"io/ioutil"
	"stockfyApi/api/handlers/fiberHandlers"
	"stockfyApi/api/middleware"
	"stockfyApi/api/presenter"
	"stockfyApi/database/postgresql"
	"stockfyApi/entity"
	"stockfyApi/usecases"
	"stockfyApi/usecases/user"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func configureSectorsApp(dbpool *pgxpool.Pool) (fiberHandlers.SectorApi,
	usecases.Applications) {

	dbInterfaces := postgresql.NewPostgresInstance(dbpool)
	firebaseInterface := user.NewExternalApi()

	applicationLogics := usecases.NewApplications(dbInterfaces,
		firebaseInterface)

	sectors := fiberHandlers.SectorApi{
		ApplicationLogic: *applicationLogics,
	}

	return sectors, *applicationLogics
}

func TestFiberHandlersIntegrationTestCreateSector(t *testing.T) {

	type body struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Error   string            `json:"error"`
		Code    int               `json:"code"`
		Sector  *presenter.Sector `json:"sector"`
	}

	type test struct {
		idToken          string
		contentType      string
		sectorName       string
		expectedResponse body
	}

	tests := []test{
		{
			idToken:     "INVALID_ID_TOKEN",
			contentType: "application/json",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Sector:  nil,
			},
		},
		{
			idToken:     "TestNoAdminID",
			contentType: "application/json",
			expectedResponse: body{
				Code:    403,
				Success: false,
				Message: entity.ErrMessageApiAuthorization.Error(),
				Error:   entity.ErrInvalidUserAdminPrivilege.Error(),
				Sector:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/pdf",
			expectedResponse: body{
				Code:    400,
				Success: false,
				Message: entity.ErrMessageApiRequest.Error(),
				Error:   entity.ErrInvalidApiBody.Error(),
				Sector:  nil,
			},
		},
		{
			idToken:     "TestAdminID",
			contentType: "application/json",
			sectorName:  "Crypto",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Sector creation was successful",
				Error:   "",
				Sector: &presenter.Sector{
					Name: "Crypto",
				},
			},
		},
	}

	DBpool := connectDatabase()

	sectors, applicationsLogics := configureSectorsApp(DBpool)

	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: applicationsLogics.UserApp,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var err error
			c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiAuthentication.Error(),
				"code":    401,
			})

			return err
		},
		ContextKey: "user",
	}))
	api.Post("/sector", sectors.CreateSector)

	for _, testCase := range tests {
		bodyResponse := body{}
		bodyRequestStruct := presenter.SectorBody{
			Sector: testCase.sectorName,
		}

		resp, _ := fiberHandlers.MockHttpRequest(app, "POST", "/api/sector",
			testCase.contentType, testCase.idToken, bodyRequestStruct)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Sector != nil {
			assert.Equal(t, testCase.expectedResponse.Sector.Name,
				bodyResponse.Sector.Name)
		} else {
			assert.Nil(t, testCase.expectedResponse.Sector)
		}

	}
}

func TestFiberHandlersIntegrationTestGetSector(t *testing.T) {

	type body struct {
		Success bool              `json:"success"`
		Message string            `json:"message"`
		Error   string            `json:"error"`
		Code    int               `json:"code"`
		Sector  *presenter.Sector `json:"sector"`
	}

	type test struct {
		idToken          string
		sectorName       string
		expectedResponse body
	}

	tests := []test{
		{
			idToken: "INVALID_ID_TOKEN",
			expectedResponse: body{
				Code:    401,
				Success: false,
				Message: entity.ErrMessageApiAuthentication.Error(),
				Error:   "",
				Sector:  nil,
			},
		},
		{
			idToken:    "TestAdminID",
			sectorName: "ERROR",
			expectedResponse: body{
				Code:    404,
				Success: false,
				Message: entity.ErrMessageApiSectorName.Error(),
				Error:   entity.ErrInvalidSectorSearchName.Error(),
				Sector:  nil,
			},
		},
		{
			idToken:    "TestAdminID",
			sectorName: "Crypto",
			expectedResponse: body{
				Code:    200,
				Success: true,
				Message: "Sector information returned successfully",
				Error:   "",
				Sector: &presenter.Sector{
					Name: "Crypto",
				},
			},
		},
	}

	DBpool := connectDatabase()

	sectors, applicationsLogics := configureSectorsApp(DBpool)

	app := fiber.New()
	api := app.Group("/api")
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: applicationsLogics.UserApp,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var err error
			c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": entity.ErrMessageApiAuthentication.Error(),
				"code":    401,
			})

			return err
		},
		ContextKey: "user",
	}))
	api.Get("/sector/:sector", sectors.GetSector)

	for _, testCase := range tests {
		bodyResponse := body{}

		resp, _ := fiberHandlers.MockHttpRequest(app, "GET", "/api/sector/"+
			testCase.sectorName, "", testCase.idToken, nil)

		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &bodyResponse)
		bodyResponse.Code = resp.StatusCode

		assert.NotNil(t, resp)
		assert.Equal(t, testCase.expectedResponse.Code, bodyResponse.Code)
		assert.Equal(t, testCase.expectedResponse.Success, bodyResponse.Success)
		assert.Equal(t, testCase.expectedResponse.Message, bodyResponse.Message)
		assert.Equal(t, testCase.expectedResponse.Error, bodyResponse.Error)

		if testCase.expectedResponse.Sector != nil {
			assert.Equal(t, testCase.expectedResponse.Sector.Name,
				bodyResponse.Sector.Name)
		} else {
			assert.Nil(t, testCase.expectedResponse.Sector)
		}

	}
}
