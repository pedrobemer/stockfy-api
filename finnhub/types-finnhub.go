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

type CompanyProfile2 struct {
	Country              string  `json:"country,omitempty"`
	Currency             string  `json:"currency,omitempty"`
	Exchange             string  `json:"exchange,omitempty"`
	Ipo                  string  `json:"ipo,omitempty"`
	MarketCapitalization float64 `json:"marketCapitalization,omitempty"`
	Name                 string  `json:"name,omitempty"`
	Phone                string  `json:"phone,omitempty"`
	ShareOutstanding     float64 `json:"shareOutstanding,omitempty"`
	Ticker               string  `json:"ticker,omitempty"`
	Weburl               string  `json:"weburl,omitempty"`
	Logo                 string  `json:"logo,omitempty"`
	FinnhubIndustry      string  `json:"finnhubIndustry,omitempty"`
}

var SymbolTypesFinnhub = map[string]string{
	"Common Stock": "STOCK",
	"ETP":          "ETF",
	"REIT":         "REIT",
}

type FinnhubApi struct {
	token string
}
