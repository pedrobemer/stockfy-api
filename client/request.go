package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func RequestAndAssignToBody(method string, url string, contentType string,
	bodyReq io.Reader, bodyResp interface{}) {
	spaceClient := http.Client{
		Timeout: time.Second * 5, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(method, url, bodyReq)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", contentType)

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &bodyResp)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
