package models

// Ticker represents market ticker data
type Ticker struct {
	Bid    float64 `json:",string"`
	Ask    float64 `json:",string"`
	Last   float64 `json:",string"`
	Symbol string
}

// Symbol represents trading pair information
type Symbol struct {
	Id                   string
	BaseCurrency         string
	QuoteCurrency        string
	FeeCurrency          string
	TakeLiquidityRate    float64 `json:",string"`
	ProvideLiquidityRate float64 `json:",string"`
	QuantityIncrement    float64 `json:",string"`
	TickSize             float64 `json:",string"`
}

// TradingSymbol represents a trading action
type TradingSymbol struct {
	Symbol   Symbol
	Action   string
	Quantity float64
	Price    float64
}

// Order represents an order on the exchange
type Order struct {
	Id            int
	ClientOrderId string
	Symbol        string
	Side          string
	Status        string
	Type          string
	TimeInForce   string
	Quantity      float64 `json:",string"`
	Price         float64 `json:",string"`
	CumQuantity   float64 `json:",string"`
	PostOnly      bool    `json:",bool"`
	CreatedAt     string
	UpdatedAt     string
}

// Balance represents account balance
type Balance struct {
	Currency  string
	Available float64 `json:",string"`
	Reserved  float64 `json:",string"`
}
