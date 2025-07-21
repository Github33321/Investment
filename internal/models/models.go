package models

import "time"

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
	ID             int       `db:"id" json:"id"`
	TotalInput     float64   `db:"total_input" json:"total_input"`
	TotalOutput    float64   `db:"total_output" json:"total_output"`
	Turnover       float64   `db:"turnover" json:"turnover"`
	TotalBuys      float64   `db:"total_buys" json:"total_buys"`
	TotalSells     float64   `db:"total_sells" json:"total_sells"`
	PortfolioValue float64   `db:"portfolio_value" json:"portfolio_value"`
	Commissions    float64   `db:"commissions" json:"commissions"`
	Taxes          float64   `db:"taxes" json:"taxes"`
	NetStockProfit float64   `db:"net_stock_profit" json:"net_stock_profit"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}
