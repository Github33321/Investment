package service

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"strings"

	"tinvest_report/internal/models"
)

func GetSummary() (model.Summary, error) {
	resp, err := http.Get("http://localhost:8082/spravka")
	if err != nil {
		return model.Summary{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Summary{}, err
	}
	var ops []model.Operation
	if err := json.Unmarshal(body, &ops); err != nil {
		return model.Summary{}, err
	}

	var totalInput, totalOutput, turnover, totalBuys, totalSells, commissions, taxes float64
	figiHoldings := make(map[string]float64)

	for _, op := range ops {
		if op.IsCanceled {
			log.Printf("[CANCELED] %s | %s | %.2f ₽", op.Date, op.OperationType, op.FloatPayment)
			continue
		}
		switch op.OperationType {
		case "OPERATION_TYPE_INPUT", "OPERATION_TYPE_INP_MULTI":
			totalInput += op.FloatPayment
		case "OPERATION_TYPE_OUTPUT", "OPERATION_TYPE_OUT_MULTI":
			totalOutput += -op.FloatPayment
		case "OPERATION_TYPE_BUY":
			if !strings.HasPrefix(op.FIGI, "FUT") {
				totalBuys += -op.FloatPayment
				turnover += -op.FloatPayment
				figiHoldings[op.FIGI] += op.Quantity
			}
		case "OPERATION_TYPE_SELL":
			if !strings.HasPrefix(op.FIGI, "FUT") {
				totalSells += op.FloatPayment
				turnover += op.FloatPayment
				figiHoldings[op.FIGI] -= op.Quantity
			}
		case "OPERATION_TYPE_BROKER_FEE", "OPERATION_TYPE_TRACK_MFEE", "OPERATION_TYPE_TRACK_PFEE":
			commissions += -op.FloatPayment
		case "OPERATION_TYPE_TAX":
			taxes += -op.FloatPayment
		}
	}

	var portfolioValue float64
	for figi, qty := range figiHoldings {
		if math.Abs(qty) < 0.0001 {
			continue
		}
		price, err := fetchPrice(figi)
		if err != nil {
			log.Printf("❌ Не удалось получить цену для %s: %v", figi, err)
			continue
		}
		portfolioValue += qty * price
	}

	netProfit := (totalSells + portfolioValue) - totalBuys - commissions - taxes

	return model.Summary{
		TotalInput:     round(totalInput),
		TotalOutput:    round(totalOutput),
		Turnover:       round(turnover),
		TotalBuys:      round(totalBuys),
		TotalSells:     round(totalSells),
		PortfolioValue: round(portfolioValue),
		Commissions:    round(commissions),
		Taxes:          round(taxes),
		NetStockProfit: round(netProfit),
	}, nil
}

func fetchPrice(figi string) (float64, error) {
	resp, err := http.Get("http://localhost:8083/figi/" + figi)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var pr model.PriceResponse
	err = json.Unmarshal(body, &pr)
	return pr.Price, err
}

func round(val float64) float64 {
	return math.Round(val*100) / 100
}
