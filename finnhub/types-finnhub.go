package finnhub

type SymbolLookupInfo struct {
	Description   string
	DisplaySymbol string
	Symbol        string
	Type          string
}

type SymbolLookupFinnhub struct {
	Count  int32
	Result []SymbolLookupInfo
}

type SymbolPriceFinnhub struct {
	C  float64
	H  float64
	L  float64
	O  float64
	PC float64
	T  float64
}

var SymbolTypesFinnhub = map[string]string{
	"Common Stock": "STOCK",
	"ETP":          "ETF",
	"REIT":         "REIT",
}
