package externalapi

import (
	"stockfyApi/externalApi/alphaVantage"
	"stockfyApi/externalApi/finnhub"
)

type ThirdPartyInterfaces struct {
	FinnhubApi      finnhub.FinnhubApi
	AlphaVantageApi alphaVantage.AlphaApi
}
