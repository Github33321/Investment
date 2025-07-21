package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"tinvest_report/internal/service"
)

type Router struct{}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.URL.Path == "/summary":
		r.handleSummary(w, req)
	case req.URL.Path == "/spravka":
		r.proxySpravka(w, req)
	case strings.HasPrefix(req.URL.Path, "/figi/"):
		r.proxyFigi(w, req)
	default:
		http.NotFound(w, req)
	}
}

func (r *Router) handleSummary(w http.ResponseWriter, req *http.Request) {
	summary, err := service.GetSummary()
	if err != nil {
		http.Error(w, "Ошибка при расчёте: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(summary)
}

func (r *Router) proxySpravka(w http.ResponseWriter, req *http.Request) {
	resp, err := http.Get("http://localhost:8082/spravka")
	if err != nil {
		http.Error(w, "Ошибка /spravka: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.Copy(w, resp.Body)
}

func (r *Router) proxyFigi(w http.ResponseWriter, req *http.Request) {
	figi := strings.TrimPrefix(req.URL.Path, "/figi/")
	if figi == "" {
		http.Error(w, "FIGI не указан", http.StatusBadRequest)
		return
	}
	resp, err := http.Get("http://localhost:8083/figi/" + figi)
	if err != nil {
		http.Error(w, "Ошибка /figi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.Copy(w, resp.Body)
}
