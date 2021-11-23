package externalapi

import (
	"stockfyApi/entity"
)

type thirdPartyInterface interface {
	VerifySymbol2(symbol string) entity.SymbolLookup
	GetPrice(symbol string) entity.SymbolPrice
	CompanyOverview(symbol string) map[string]string
}

type ThirdPartyInterfaces struct {
	FinnhubApi      thirdPartyInterface
	AlphaVantageApi thirdPartyInterface
}
