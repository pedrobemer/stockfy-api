package commonTypes

type SymbolLookup struct {
	Fullname string `json:",omitempty"`
	Symbol   string `json:",omitempty"`
	Type     string `json:",omitempty"`
}

type SymbolPrice struct {
	Symbol         string  `json:",omitempty"`
	CurrentPrice   float64 `json:",omitempty"`
	HighPrice      float64 `json:",omitempty"`
	LowPrice       float64 `json:",omitempty"`
	OpenPrice      float64 `json:",omitempty"`
	PrevClosePrice float64 `json:",omitempty"`
	MarketCap      float64 `json:",omitempty"`
}
