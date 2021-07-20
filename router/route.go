package router

import (
	"stockfyApi/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App, dbpool *pgxpool.Pool) {
	// Middleware
	api := app.Group("/api")

	// routes

	// Intermediary REST API for the Finnhub API
	api.Get("/finnhub/symbol-lookup", handlers.GetSymbolFinnhub)
	api.Get("/finnhub/symbol-price", handlers.GetSymbolPriceFinnhub)

	// Intermediary REST API for the Alpha Vantage API
	api.Get("/alpha-vantage/symbol-lookup", handlers.GetSymbolAlphaVantage)
	api.Get("/alpha-vantage/symbol-price", handlers.GetSymbolPriceAlphaVantage)

	// REST API for the assets table
	api.Get("/asset/asset-types", handlers.GetAssetsFromAssetType)
	api.Get("/asset/:symbol", handlers.GetAsset)
	api.Get("/asset/:symbol/orders", handlers.GetAssetWithOrders)
	api.Post("/asset", handlers.PostAsset)

	// REST API for the asset types table
	api.Get("/asset-types", handlers.GetAssetTypes)

	// REST API to for the sector table
	api.Get("/sector", handlers.GetAllSectors)
	api.Get("/sector/:sector", handlers.GetSector)
	api.Post("/sector", handlers.PostSector)

	// REST API for the orders table
	api.Post("/orders", handlers.PostOrder)

	// REST API for the brokerage table
	api.Get("/brokerage/:name", handlers.GetBrokerageFirm)
	api.Get("/brokerage", handlers.GetBrokerageFirms)

}
