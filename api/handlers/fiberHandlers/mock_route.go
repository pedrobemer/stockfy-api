package fiberHandlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func MockHttpRequest(app *fiber.App, method string, path string, contentType string,
	authToken string, jsonRequest interface{}) (*http.Response, error) {

	var err error
	var resp *http.Response

	bearerToken := "Bearer " + authToken

	bodyByte, err := json.Marshal(jsonRequest)

	req := httptest.NewRequest(method, path, bytes.NewReader(bodyByte))
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", bearerToken)

	resp, err = app.Test(req)

	return resp, err
}

func (m *MockClient) MockHttpOutsideRequest(method string, path string,
	contentType string, body io.Reader) (*http.Response, error) {

	req := httptest.NewRequest(method, path, body)
	req.Header.Set("Content-Type", contentType)

	resp, _ := m.Do(req)

	return resp, nil
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
