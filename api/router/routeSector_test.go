package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"stockfyApi/database"
	"stockfyApi/handlers"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

type respSector struct {
	Message string                     `json:"message,omitempty"`
	Sector  []database.SectorApiReturn `json:"sector,omitempty"`
	Success bool                       `json:"success,omitempty"`
	Error   string                     `json:"error,omitempty"`
}

func MockDatabase(query string, columns []string) (mock pgxmock.PgxConnIface,
	rows *pgxmock.Rows, err error) {
	mock, err = pgxmock.NewConn()
	if err != nil {
		return mock, rows, err
	}
	defer mock.Close(context.Background())

	rows = mock.NewRows(columns)

	return mock, rows, err
}

func MockHttpRequest(app *fiber.App, method string, path string,
	jsonResponse interface{}, jsonRequest interface{}) (*http.Response, error) {

	var err error
	var resp *http.Response

	bodyByte, err := json.Marshal(jsonRequest)

	req := httptest.NewRequest(method, path, bytes.NewReader(bodyByte))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		return resp, err
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return resp, readErr
	}

	jsonErr := json.Unmarshal(body, jsonResponse)
	if jsonErr != nil {
		return resp, jsonErr
	}

	return resp, err
}

func TestApiSectorRootValidResponse(t *testing.T) {
	var jsonResponse respSector

	listSectors := []database.SectorApiReturn{
		{
			Id:   "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Name: "Finance",
		},
		{
			Id:   "62d4d8e2-95e5-4144-b17b-0d147c98d85c",
			Name: "Technology",
		},
	}
	expectedJsonResponse := respSector{
		Message: "All sectors returned successfully",
		Sector:  listSectors,
		Success: true,
	}

	// Mock postgreSQL for API request
	columns := []string{"id", "name"}

	mock, rows, err := MockDatabase(".*", columns)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(".*").WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "Finance").
			AddRow("62d4d8e2-95e5-4144-b17b-0d147c98d85c", "Technology"))

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")

	sector := handlers.SectorApi{Db: mock}

	api.Get("/sector", sector.GetAllSectors)

	resp, err := MockHttpRequest(app, "GET", "/api/sector", &jsonResponse, nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	assert.NotNil(t, resp, "Request Not Nil")
	assert.Equal(t, 200, resp.StatusCode, "Request Status Code")
	assert.Equal(t, expectedJsonResponse, jsonResponse)

}

func TestApiSectorSingleValidResponse(t *testing.T) {
	var jsonResponse respSector

	listSectors := []database.SectorApiReturn{
		{
			Id:   "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Name: "Finance",
		},
	}
	expectedJsonResponse := respSector{
		Message: "Sector information returned successfully",
		Sector:  listSectors,
		Success: true,
	}

	// Mock postgreSQL for API request
	columns := []string{"id", "name"}

	mock, rows, err := MockDatabase(".*", columns)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(".*").WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "Finance"))

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")

	sector := handlers.SectorApi{Db: mock}

	api.Get("/sector/:sector", sector.GetSector)

	resp, err := MockHttpRequest(app, "GET", "/api/sector/Finance",
		&jsonResponse, nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	assert.NotNil(t, resp, "Request Not Nil")
	assert.Equal(t, 200, resp.StatusCode, "Request Status Code")
	assert.Equal(t, expectedJsonResponse, jsonResponse)
}

func TestApiSectorSingleInvalidResponse(t *testing.T) {
	var jsonResponse respSector
	// var jsonResponse interface{}

	expectedJsonResponse := respSector{
		Error:   "FetchSector: Nonexistent sector in the database",
		Success: false,
	}

	// Mock postgreSQL for API request
	columns := []string{"id", "name"}

	mock, _, err := MockDatabase(".*", columns)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(".*").WillReturnError(fmt.Errorf("some error"))

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")

	sector := handlers.SectorApi{Db: mock}

	api.Get("/sector/:sector", sector.GetSector)

	resp, err := MockHttpRequest(app, "GET", "/api/sector/Finance",
		&jsonResponse, nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	assert.NotNil(t, resp, "Request Not Nil")
	assert.Equal(t, 500, resp.StatusCode, "Request Status Code")
	assert.Equal(t, expectedJsonResponse, jsonResponse)
}

func TestApiSectorSingleUnauthorized(t *testing.T) {
	var jsonResponse respSector

	expectedJsonResponse := respSector{
		Error:   "Unauthorized Sector Search",
		Success: false,
	}

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")

	sector := handlers.SectorApi{}

	api.Get("/sector/:sector", sector.GetSector)

	resp, err := MockHttpRequest(app, "GET", "/api/sector/ALL", &jsonResponse,
		nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	assert.NotNil(t, resp, "Request Not Nil")
	assert.Equal(t, 500, resp.StatusCode, "Request Status Code")
	assert.Equal(t, expectedJsonResponse, jsonResponse)
}

func TestApiSectorPostSector(t *testing.T) {
	var jsonResponse respSector

	jsonRequest := database.SectorBodyPost{
		Sector: "Finance",
	}

	listSectors := []database.SectorApiReturn{
		{
			Id:   "0a52d206-ed8b-11eb-9a03-0242ac130003",
			Name: "Finance",
		},
	}

	expectedJsonResponse := respSector{
		Message: "Created sector successfully",
		Success: true,
		Sector:  listSectors,
	}

	// Mock postgreSQL for API request
	columns := []string{"id", "name"}
	mock, rows, err := MockDatabase(".*", columns)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(".*").WillReturnRows(
		rows.AddRow("0a52d206-ed8b-11eb-9a03-0242ac130003", "Finance"))

	// Mock HTTP request
	app := fiber.New()
	api := app.Group("/api")

	sector := handlers.SectorApi{Db: mock}

	api.Post("/sector", sector.PostSector)

	resp, err := MockHttpRequest(app, "POST", "/api/sector", &jsonResponse,
		jsonRequest)
	if err != nil {
		t.Fatalf("%s", err)
	}

	assert.NotNil(t, resp, "Request Not Nil")
	assert.Equal(t, 200, resp.StatusCode, "Request Status Code")
	assert.Equal(t, expectedJsonResponse, jsonResponse)
}
