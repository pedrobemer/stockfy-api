package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func MockHttpRequest(app *fiber.App, method string, path string, authToken string,
	jsonRequest interface{}) (*http.Response, error) {

	var err error
	var resp *http.Response

	bearerToken := "Bearer " + authToken

	bodyByte, err := json.Marshal(jsonRequest)

	req := httptest.NewRequest(method, path, bytes.NewReader(bodyByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", bearerToken)

	resp, err = app.Test(req)
	if err != nil {
		return resp, err
	}

	return resp, err
}
