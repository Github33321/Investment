package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"tinvest_report/internal/models"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) SaveSummary(ctx context.Context, summary models.Summary) error {
	query := `
	INSERT INTO summary (
		total_input, total_output, turnover, total_buys,
		total_sells, portfolio_value, commissions, taxes, net_stock_profit, created_at
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9, now())`

	_, err := r.DB.Exec(ctx, query,
		summary.TotalInput, summary.TotalOutput, summary.Turnover, summary.TotalBuys,
		summary.TotalSells, summary.PortfolioValue, summary.Commissions,
		summary.Taxes, summary.NetStockProfit,
	)
	return err
}
func (r *Repository) GetSummaries(ctx context.Context, date string) ([]models.Summary, error) {
	var summaries []models.Summary

	if date != "" {
		query := `
			SELECT * FROM summary
			WHERE created_at >= $1 AND created_at < $2
			ORDER BY created_at DESC
		`
		start, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		end := start.Add(24 * time.Hour)

		rows, err := r.DB.Query(ctx, query, start, end)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var s models.Summary
			err := rows.Scan(
				&s.ID, &s.TotalInput, &s.TotalOutput, &s.Turnover, &s.TotalBuys,
				&s.TotalSells, &s.PortfolioValue, &s.Commissions, &s.Taxes,
				&s.NetStockProfit, &s.CreatedAt,
			)
			if err != nil {
				return nil, err
			}
			summaries = append(summaries, s)
		}
		return summaries, nil
	}

	// Без фильтра по дате
	query := `SELECT * FROM summary ORDER BY created_at DESC`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Summary
		err := rows.Scan(
			&s.ID, &s.TotalInput, &s.TotalOutput, &s.Turnover, &s.TotalBuys,
			&s.TotalSells, &s.PortfolioValue, &s.Commissions, &s.Taxes,
			&s.NetStockProfit, &s.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}
