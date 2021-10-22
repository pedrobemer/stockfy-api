package usecases

import (
	"stockfyApi/usecases/asset"
	dbverification "stockfyApi/usecases/dbVerification"
	"stockfyApi/usecases/sector"
	"stockfyApi/usecases/user"
)

// type MockApplications struct {
// 	// AssetApp          asset.MockApplication
// 	// AssetTypeApp      assettype.MockApplication
// 	// AssetUserApp      assetusers.MockApplication
// 	SectorApp sector.MockApplication
// 	// UserApp           user.MockApplication
// 	// OrderApp          order.MockApplication
// 	// BrokerageApp      brokerage.MockApplication
// 	// EarningsApp       earnings.MockApplication
// 	// DbVerificationApp dbverification.MockApplication
// }

func NewMockApplications() *Applications {
	return &Applications{
		SectorApp: sector.NewMockApplication(),
		// AssetTypeApp:      *assettype.NewApplication(),
		AssetApp: asset.NewMockApplication(),
		// AssetUserApp:      *assetusers.NewApplication(),
		UserApp: user.NewMockApplication(),
		// OrderApp:          *order.NewApplication(),
		// BrokerageApp:      *brokerage.NewApplication(),
		// EarningsApp:       *earnings.NewApplication(),
		DbVerificationApp: dbverification.NewMockApplication(),
	}
}
