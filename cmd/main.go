package main

import (
	"log"
	"net/http"

	"tinvest_report/internal/handlers"
)

func main() {
	router := handlers.NewRouter()
	log.Println("✅ Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
