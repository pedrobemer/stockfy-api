package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/actionName", func(c *fiber.Ctx) error {
		payload := Payload{}

		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		Symbol := payload.Input.Arg1.Symbol
		url := "https://finnhub.io/api/v1/quote?symbol=" + Symbol + "&token=c2o3062ad3ie71thpra0"

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

		apple := stock{}
		jsonErr := json.Unmarshal(body, &apple)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		return c.JSON(apple)
	})

	app.Listen(":3000")
}
