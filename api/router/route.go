package router

import (
	"log"
	"stockfyApi/api/handlers/fiberHandlers"
	"stockfyApi/api/middleware"
	externalapi "stockfyApi/externalApi"
	"stockfyApi/usecases"
	"stockfyApi/usecases/logicApi"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(chosenApi string, firebaseKey string,
	usecases *usecases.Applications,
	externalInterfaces externalapi.ThirdPartyInterfaces) {

	switch chosenApi {
	case "FIBER":
		fiberRoutes(firebaseKey, usecases, externalInterfaces)
		break
	default:
		log.Panic("Wrong chosen API. Only Fiber is available.")
	}

}

func fiberRoutes(firebaseKey string, usecases *usecases.Applications,
	externalInterfaces externalapi.ThirdPartyInterfaces) {
	app := fiber.New()

	logicApiUseCases := logicApi.NewApplication(*usecases, externalInterfaces)

	// REST API Handlers
	// sector := fiberHandlers.SectorApi{ApplicationLogic: *usecases}
	asset := fiberHandlers.AssetApi{
		ApplicationLogic:   *usecases,
		ExternalInterfaces: externalInterfaces,
		LogicApi:           *logicApiUseCases,
	}
	// assetType := fiberHandlers.AssetTypeApi{ApplicationLogic: *usecases}
	order := fiberHandlers.OrderApi{
		ApplicationLogic:   *usecases,
		ExternalInterfaces: externalInterfaces,
		LogicApi:           *logicApiUseCases,
	}
	brokerage := fiberHandlers.BrokerageApi{
		ApplicationLogic: *usecases,
	}
	earnings := fiberHandlers.EarningsApi{
		ApplicationLogic: *usecases,
		ApiLogic:         *logicApiUseCases,
	}
	alpha := fiberHandlers.AlphaVantageApi{
		ApplicationLogic: *usecases,
		Api:              &externalInterfaces.AlphaVantageApi,
	}
	finn := fiberHandlers.FinnhubApi{
		ApplicationLogic: *usecases,
		Api:              &externalInterfaces.FinnhubApi,
	}
	firebaseApi := fiberHandlers.FirebaseApi{
		ApplicationLogic: *usecases,
		FirebaseWebKey:   firebaseKey,
	}

	api := app.Group("/api")

	// REST API to create a user on Firebase
	api.Post("/signup", firebaseApi.SignUp)
	api.Post("/forgot-password", firebaseApi.ForgotPassword)

	// Middleware
	api.Use(middleware.NewFiberMiddleware(middleware.FiberMiddleware{
		UserAuthentication: &usecases.UserApp,
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
	api.Get("/finnhub/symbol-lookup", finn.GetSymbol)
	api.Get("/finnhub/symbol-price", finn.GetSymbolPrice)
	// api.Get("/finnhub/company-profile", finn.GetCompanyProfile2Finnhub)

	// Intermediary REST API for the Alpha Vantage API
	api.Get("/alpha-vantage/symbol-lookup", alpha.GetSymbol)
	api.Get("/alpha-vantage/symbol-price", alpha.GetSymbolPrice)
	// api.Get("/alpha-vantage/company-overview", alpha.GetCompanyOverviewAlphaVantage)

	// REST API for the assets table
	api.Get("/asset/asset-types", asset.GetAssetsFromAssetType)
	api.Get("/asset/:symbol", asset.GetAsset)
	api.Get("/asset/:symbol/orders", asset.GetAssetWithOrders)
	api.Post("/asset", asset.CreateAsset)
	api.Delete("/asset/:symbol", asset.DeleteAsset)

	// REST API for the asset types table
	// api.Get("/asset-types", assetType.GetAssetTypes)

	// REST API to for the sector table
	// api.Get("/sector", sector.GetAllSectors)
	// api.Get("/sector/:sector", sector.GetSector)
	// api.Post("/sector", sector.PostSector)

	// REST API for the orders table
	api.Get("/orders", order.GetOrdersFromAssetUser)
	api.Post("/orders", order.CreateUserOrder)
	api.Delete("orders/:id", order.DeleteOrderFromUser)
	api.Put("/orders/:id", order.UpdateOrderFromUser)

	// REST API for the brokerage table
	// api.Get("/brokerage/:name", brokerage.GetBrokerageFirm)
	api.Get("/brokerage", brokerage.GetBrokerageFirms)

	// REST API for the earning table
	api.Get("/earnings", earnings.GetEarningsFromAssetUser)
	api.Post("/earnings", earnings.CreateEarnings)
	api.Put("/earnings/:id", earnings.UpdateEarningFromUser)
	api.Delete("/earnings/:id", earnings.DeleteEarningFromUser)

	app.Listen(":3000")

}
