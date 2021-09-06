package main

import (
	"context"
	"fmt"
	"os"
	"stockfyApi/database"
	"stockfyApi/router"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	DB_USER := viperReadEnvVariable("DB_USER")
	DB_PASSWORD := viperReadEnvVariable("DB_PASSWORD")
	DB_NAME := viperReadEnvVariable("DB_NAME")
	FIREBASE_API_WEB_KEY := viperReadEnvVariable("FIREBASE_API_WEB_KEY")

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	database.DBpool, err = pgx.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer database.DBpool.Close(context.Background())

	app := fiber.New()

	router.SetupRoutes(app, FIREBASE_API_WEB_KEY)

	app.Listen(":3000")
}
