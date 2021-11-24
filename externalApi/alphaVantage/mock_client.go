package alphaVantage

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"stockfyApi/api/handlers/fiberHandlers"
)

type MockClient struct {
	Client fiberHandlers.MockClient
}

func (mc *MockClient) HttpOutsideClientRequest(method string, url string,
	contentType string, bodyReq io.Reader, bodyResp interface{}) {

	resp, _ := mc.Client.MockHttpOutsideRequest(method, url, contentType,
		bodyReq)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &bodyResp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
