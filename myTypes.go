package main

type SymbolPriceNotFormatted struct {
	C  float64
	H  float64
	L  float64
	O  float64
	PC float64
	T  float64
}

type SymbolPrice struct {
	CurrentPrice   float64
	HighPrice      float64
	LowPrice       float64
	OpenPrice      float64
	PrevClosePrice float64
	MarketCap      float64
}

type SymbolLookupInfo struct {
	Description   string
	DisplaySymbol string
	Symbol        string
	Type          string
}

type SymbolLookup struct {
	Count  int32
	Result []SymbolLookupInfo
}
