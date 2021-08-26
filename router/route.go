package router

import (
	"stockfyApi/database"
	"stockfyApi/handlers"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api")

	// Handlers
	sector := handlers.SectorApi{Db: database.DBpool}
	asset := handlers.AssetApi{Db: database.DBpool}
	assetType := handlers.AssetTypeApi{Db: database.DBpool}
	order := handlers.OrderApi{Db: database.DBpool}
	brokerage := handlers.BrokerageApi{Db: database.DBpool}
	earnings := handlers.EarningsApi{Db: database.DBpool}
	alpha := handlers.AlphaVantageApi{}
	finn := handlers.FinnhubApi{}

	// Intermediary REST API for the Finnhub API
	api.Get("/finnhub/symbol-lookup", finn.GetSymbolFinnhub)
	api.Get("/finnhub/symbol-price", finn.GetSymbolPriceFinnhub)
	api.Get("/finnhub/company-profile", finn.GetCompanyProfile2Finnhub)

	// Intermediary REST API for the Alpha Vantage API
	api.Get("/alpha-vantage/symbol-lookup", alpha.GetSymbolAlphaVantage)
	api.Get("/alpha-vantage/symbol-price", alpha.GetSymbolPriceAlphaVantage)
	api.Get("/alpha-vantage/company-overview", alpha.GetCompanyOverviewAlphaVantage)

	// REST API for the assets table
	api.Get("/asset/asset-types", asset.GetAssetsFromAssetType)
	api.Get("/asset/:symbol", asset.GetAsset)
	api.Get("/asset/:symbol/orders", asset.GetAssetWithOrders)
	api.Post("/asset", asset.PostAsset)
	api.Delete("/asset/:symbol", asset.DeleteAsset)

	// REST API for the asset types table
	api.Get("/asset-types", assetType.GetAssetTypes)

	// REST API to for the sector table
	api.Get("/sector", sector.GetAllSectors)
	api.Get("/sector/:sector", sector.GetSector)
	api.Post("/sector", sector.PostSector)

	// REST API for the orders table
	api.Post("/orders", order.PostOrder)
	api.Delete("orders/:id", order.DeleteOrder)
	api.Put("/orders/:id", order.UpdateOrder)

	// REST API for the brokerage table
	api.Get("/brokerage/:name", brokerage.GetBrokerageFirm)
	api.Get("/brokerage", brokerage.GetBrokerageFirms)

	// REST API for the earning table
	api.Post("/earnings", earnings.PostEarnings)

}
