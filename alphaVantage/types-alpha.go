package alphaVantage

type SymbolLookupAlpha struct {
	BestMatches []SymbolLookupInfo `json:"bestMatches"`
}

type SymbolLookupInfo struct {
	Symbol      string `json:"1. symbol"`
	Name        string `json:"2. name"`
	Type        string `json:"3. type"`
	Region      string `json:"4. region"`
	MarketOpen  string `json:"5. marketOpen"`
	MarketClose string `json:"6. marketClose"`
	Timezone    string `json:"7. timezone"`
	Currency    string `json:"8. currency"`
	MatchScore  string `json:"9. matchScore"`
}

type SymbolPriceAlpha struct {
	GlobalQuote SymbolPriceInfo `json:"Global Quote"`
}

type SymbolPriceInfo struct {
	Symbol        string `json:"01. symbol"`
	Open          string `json:"02. open"`
	High          string `json:"03. high"`
	Low           string `json:"04. low"`
	Price         string `json:"05. price"`
	Volume        string `json:"06. volume"`
	LatestDay     string `json:"07. latest trading day"`
	PrevClose     string `json:"08. previous close"`
	Change        string `json:"09. change"`
	ChangePercent string `json:"10. change percent"`
}

type CompanyOverview map[string]string

var ListValidBrETF = [5]string{"BOVA11", "SMAL11", "IVVB11", "HASH11", "ECOO11"}

type AlphaApi struct {
	Token string
}
