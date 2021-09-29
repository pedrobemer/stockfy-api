package router

import (
	"stockfyApi/api/handlers/fiberHandlers"
	"stockfyApi/database"
	"stockfyApi/externalApi/alphaVantage"
	"stockfyApi/usecases"

	// "stockfyApi/fiberHandlers"
	"stockfyApi/api/middleware"
	"stockfyApi/externalApi/firebaseApi"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, firebaseKey string, usecases *usecases.Applications) {

	auth := firebaseApi.SetupFirebase(
		"stockfy-api-firebase-adminsdk-cwuka-f2c828fb90.json")

	// REST API Handlers
	sector := fiberHandlers.SectorApi{ApplicationLogic: *usecases}
	asset := fiberHandlers.AssetApi{ApplicationLogic: *usecases}
	assetType := fiberHandlers.AssetTypeApi{ApplicationLogic: *usecases}
	order := fiberHandlers.OrderApi{ApplicationLogic: *usecases}
	brokerage := fiberHandlers.BrokerageApi{ApplicationLogic: *usecases}
	earnings := fiberHandlers.EarningsApi{ApplicationLogic: *usecases}
	alpha := fiberHandlers.AlphaVantageApi{
		Api: alphaVantage.AlphaApi{Token: firebaseKey},
	}
	finn := fiberHandlers.FinnhubApi{}
	firebaseApi := fiberHandlers.FirebaseApi{Db: database.DBpool, FirebaseAuth: auth,
		FirebaseWebKey: firebaseKey}

	// Middleware
	api := app.Group("/api")

	// REST API to create a user on Firebase
	api.Post("/signup", firebaseApi.SignUp)
	api.Post("/forgot-password", firebaseApi.ForgotPassword)

	api.Use(middleware.NewFirebase(middleware.Firebase{
		FirebaseAuth: auth,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			var err error
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "idToken unauthorized",
			})

			return err
		},
		ContextKey: "user",
	}))

	// REST API to disable, delete and update User information
	api.Post("/delete-user", firebaseApi.DeleteUser)
	api.Post("/update-user", firebaseApi.UpdateUserInfo)

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
	api.Post("/asset", asset.CreateAsset)
	api.Delete("/asset/:symbol", asset.DeleteAsset)

	// REST API for the asset types table
	api.Get("/asset-types", assetType.GetAssetTypes)

	// REST API to for the sector table
	api.Get("/sector", sector.GetAllSectors)
	api.Get("/sector/:sector", sector.GetSector)
	api.Post("/sector", sector.PostSector)

	// REST API for the orders table
	api.Get("/orders", order.GetOrdersFromAssetUser)
	api.Post("/orders", order.PostOrderFromUser)
	api.Delete("orders/:id", order.DeleteOrderFromUser)
	api.Put("/orders/:id", order.UpdateOrderFromUser)

	// REST API for the brokerage table
	api.Get("/brokerage/:name", brokerage.GetBrokerageFirm)
	api.Get("/brokerage", brokerage.GetBrokerageFirms)

	// REST API for the earning table
	api.Get("/earnings", earnings.GetEarningsFromAssetUser)
	api.Post("/earnings", earnings.PostEarnings)
	api.Put("/earnings/:id", earnings.UpdateEarningFromUser)
	api.Delete("/earnings/:id", earnings.DeleteEarningFromUser)

}
