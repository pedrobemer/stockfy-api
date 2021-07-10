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
	Symbol         string
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

type SymbolLookupFinnhub struct {
	Count  int32
	Result []SymbolLookupInfo
}

type SymbolLookupAlpha struct {
	Count  int32
	Result []SymbolLookupInfo
}

type Resolver struct{}

// FinhubArgs is the args to finhub
type FinhubArgs struct {
	Symbol string
}

type SymbolList struct {
	Symbols []string
}
