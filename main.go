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

const (
	DB_USER     = "pedrobemer"
	DB_PASSWORD = "pirulito"
	DB_NAME     = "stockfy"
)

func main() {
	var err error

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	database.DBpool, err = pgx.Connect(context.Background(), dbinfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer database.DBpool.Close(context.Background())

	app := fiber.New()

	router.SetupRoutes(app)

	app.Listen(":3000")
}
