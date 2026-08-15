package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	router := chi.NewRouter()
	router.Get("/health", healthHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("APIサーバーをポート8080で起動します")
	log.Fatal(server.ListenAndServe())
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
		log.Printf("ヘルスチェック応答の書き込みに失敗しました: %v", err)
	}
}
