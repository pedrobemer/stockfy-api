package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func RequestAndAssignToBody(url string, anyThing interface{}) {
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

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

	jsonErr := json.Unmarshal(body, &anyThing)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
