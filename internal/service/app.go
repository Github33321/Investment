package service

import (
	"tinvest_report/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Tinkoff *TinkoffClient
	Repo    *repository.Repository
}

func NewApp(db *pgxpool.Pool) *App {
	return &App{
		Tinkoff: NewTinkoffClient(),
		Repo:    repository.NewRepository(db),
	}
}
