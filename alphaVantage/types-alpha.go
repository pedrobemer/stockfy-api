package alphaVantage

type SymbolLookupAlpha struct {
	BestMatches []map[string]string
}

type SymbolLookupInfo map[string]string

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

var ListValidBrETF = [5]string{"BOVA11", "SMAL11", "IVVB11", "HASH11", "ECOO11"}
