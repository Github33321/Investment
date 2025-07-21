package main

import (
	"log"
	"net/http"

	"tinvest_report/internal/handlers"
	"tinvest_report/internal/service"
)

func main() {
	app := service.NewApp()

	handler := handlers.NewHandler(app)

	http.HandleFunc("/summary", handler.SummaryHandler)
	http.HandleFunc("/spravka", handler.SpravkaHandler)
	http.HandleFunc("/figi/", handler.FigiHandler)

	log.Println("✅ Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
