package models

type StockResponse struct {
	Price     float64 `json:"c"`
	PriceHigh float64 `json:"h"`
	PriceLow  float64 `json:"l"`
}
