// @title TInvest Report API
// @version 1.0
// @description API для генерации и получения отчётов по операциям Tinkoff Invest.
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	_ "tinvest_report/docs"

	"tinvest_report/db"
	"tinvest_report/internal/handlers"
	"tinvest_report/internal/service"
	"tinvest_report/internal/tasks"
)

func main() {
	err := godotenv.Load(filepath.Join(".", ".env"))
	if err != nil {
		log.Println("Не удалось загрузить .env, используется значение по умолчанию")
	}
	dsn := os.Getenv("POSTGRES_DSN")

	pool, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}

	app := service.NewApp(pool)
	handler := handlers.NewHandler(app)
	http.Handle("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/summary", handler.SummaryHandler)
	http.HandleFunc("/spravka", handler.SpravkaHandler)
	http.HandleFunc("/figi/", handler.FigiHandler)
	http.HandleFunc("/summary/save", handler.SaveSummaryHandler)
	http.HandleFunc("/summaries", handler.GetSummariesHandler)

	tasks.AutoSaveSummary(1 * time.Hour)

	log.Println("✅ Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
