package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"

	"tinvest_report/internal/service"
)

type Handler struct {
	app *service.App
}

func NewHandler(app *service.App) *Handler {
	return &Handler{app: app}
}

func (h *Handler) SpravkaHandler(w http.ResponseWriter, r *http.Request) {
	ops, err := h.app.Tinkoff.GetOperations()
	if err != nil {
		http.Error(w, "Ошибка получения операций: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(ops)
}

func (h *Handler) FigiHandler(w http.ResponseWriter, r *http.Request) {
	figi := strings.TrimPrefix(r.URL.Path, "/figi/")
	if figi == "" {
		http.Error(w, "FIGI не указан", http.StatusBadRequest)
		return
	}
	priceData, err := h.app.Tinkoff.GetFigiPrice(figi)
	if err != nil {
		http.Error(w, "Ошибка получения цены: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(priceData)
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

func (h *Handler) SummaryHandler(w http.ResponseWriter, r *http.Request) {
	ops, err := h.app.Tinkoff.GetOperations()
	if err != nil {
		http.Error(w, "Ошибка загрузки операций: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var totalInput, totalOutput, turnover, totalBuys, totalSells, commissions, taxes float64
	figiHoldings := make(map[string]float64)

	log.Println("🔍 Подробности по операциям:")
	for _, op := range ops {
		if op.IsCanceled {
			log.Printf("[CANCELED] %s | %s | %.2f ₽", op.Date, op.Operation, op.FloatPayment)
			continue
		}

		switch op.Operation {
		case "OPERATION_TYPE_INPUT", "OPERATION_TYPE_INP_MULTI":
			totalInput += op.FloatPayment
		case "OPERATION_TYPE_OUTPUT", "OPERATION_TYPE_OUT_MULTI":
			totalOutput += -op.FloatPayment
		case "OPERATION_TYPE_BUY":
			if !strings.HasPrefix(op.Figi, "FUT") {
				totalBuys += -op.FloatPayment
				turnover += -op.FloatPayment
				figiHoldings[op.Figi] += op.Quantity
			}
		case "OPERATION_TYPE_SELL":
			if !strings.HasPrefix(op.Figi, "FUT") {
				totalSells += op.FloatPayment
				turnover += op.FloatPayment
				figiHoldings[op.Figi] -= op.Quantity
			}
		case "OPERATION_TYPE_BROKER_FEE", "OPERATION_TYPE_TRACK_MFEE", "OPERATION_TYPE_TRACK_PFEE":
			commissions += -op.FloatPayment
		case "OPERATION_TYPE_TAX":
			taxes += -op.FloatPayment
		}
	}

	var portfolioValue float64
	for figi, qty := range figiHoldings {
		if figi == "" || math.Abs(qty) < 0.0001 {
			continue
		}
		if math.Abs(qty) < 0.0001 {
			continue
		}
		priceData, err := h.app.Tinkoff.GetFigiPrice(figi)
		if err != nil {
			log.Printf("❌ Не удалось получить цену для %s: %v", figi, err)
			continue
		}
		portfolioValue += qty * priceData.Price
	}

	netProfit := (totalSells + portfolioValue) - totalBuys - commissions - taxes

	summary := Summary{
		TotalInput:     math.Round(totalInput*100) / 100,
		TotalOutput:    math.Round(totalOutput*100) / 100,
		Turnover:       math.Round(turnover*100) / 100,
		TotalBuys:      math.Round(totalBuys*100) / 100,
		TotalSells:     math.Round(totalSells*100) / 100,
		PortfolioValue: math.Round(portfolioValue*100) / 100,
		Commissions:    math.Round(commissions*100) / 100,
		Taxes:          math.Round(taxes*100) / 100,
		NetStockProfit: math.Round(netProfit*100) / 100,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		http.Error(w, "Ошибка кодирования", http.StatusInternalServerError)
	}
}
