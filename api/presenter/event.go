package presenter

type EventBody struct {
	Symbol         string  `json:"symbol"`
	SymbolDemerger string  `json:"symbolDemerger"`
	EventRate      float64 `json:"eventRate"`
	Price          float64 `json:"price"`
	Currency       string  `json:"currency"`
	EventType      string  `json:"eventType"`
	Date           string  `json:"date"`
}
