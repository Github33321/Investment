package tasks

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"tinvest_report/internal/models"
)

func AutoSaveSummary(interval time.Duration) {
	go func() {
		for {
			log.Println("⏱ Автосохранение summary...")
			saveSummaryOnce()
			time.Sleep(interval)
		}
	}()
}

func saveSummaryOnce() {
	resp, err := http.Get("http://localhost:8080/summary")
	if err != nil {
		log.Println("⚠️ Ошибка запроса /summary:", err)
		return
	}
	defer resp.Body.Close()

	var summary models.Summary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		log.Println("⚠️ Ошибка декодирования summary:", err)
		return
	}

	body, err := json.Marshal(summary)
	if err != nil {
		log.Println("⚠️ Ошибка маршалинга summary:", err)
		return
	}

	_, err = http.Post("http://localhost:8080/summary/save", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("⚠️ Ошибка POST /summary/save:", err)
		return
	}

	log.Println("✅ Summary сохранен")
}
