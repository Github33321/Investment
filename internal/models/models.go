package model

type Operation struct {
	ID            string  `json:"id"`
	Currency      string  `json:"currency"`
	FloatPayment  float64 `json:"float_payment"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	OperationType string  `json:"operation_type"`
	FIGI          string  `json:"figi"`
	Quantity      float64 `json:"quantity"`
	Price         float64 `json:"price"`
	IsCanceled    bool    `json:"is_canceled"`
}

type PriceResponse struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Summary struct {
	TotalInput     float64 `json:"total_input"`
	TotalOutput    float64 `json:"total_output"`
	Turnover       float64 `json:"turnover"`
	TotalBuys      float64 `json:"total_buys"`
	TotalSells     float64 `json:"total_sells"`
	PortfolioValue float64 `json:"portfolio_value"`
	Commissions    float64 `json:"commissions"`
	Taxes          float64 `json:"taxes"`
	NetStockProfit float64 `json:"net_stock_profit"`
}
