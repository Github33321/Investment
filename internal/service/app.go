package service

type App struct {
	Tinkoff *TinkoffClient
}

func NewApp() *App {
	return &App{
		Tinkoff: NewTinkoffClient(),
	}
}
