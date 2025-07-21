package client

import (
	"encoding/json"
	"io"
	"net/http"

	"tinvest_report/internal/models"
)

func FetchOperations() ([]models.Operation, error) {
	resp, err := http.Get("http://localhost:8080/spravka")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ops []models.Operation
	err = json.Unmarshal(body, &ops)
	return ops, err
}

func FetchPrice(figi string) (float64, error) {
	resp, err := http.Get("http://localhost:8080/figi/" + figi)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var pr models.PriceResponse
	err = json.Unmarshal(body, &pr)
	return pr.Price, err
}
