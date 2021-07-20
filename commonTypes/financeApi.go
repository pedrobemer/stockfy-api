package commonTypes

type SymbolLookup struct {
	Fullname string
	Symbol   string
	Type     string
}

type SymbolPrice struct {
	Symbol         string
	CurrentPrice   float64
	HighPrice      float64
	LowPrice       float64
	OpenPrice      float64
	PrevClosePrice float64
	MarketCap      float64
}
