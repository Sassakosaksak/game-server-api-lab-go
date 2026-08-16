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
	server := &http.Server{
		Addr:    ":8080",
		Handler: newRouter(),
	}

	log.Println("APIサーバーをポート8080で起動します")
	log.Fatal(server.ListenAndServe())
}

// 本番起動とテストで同じルーティング設定を使うためにRouterを作る。
func newRouter() http.Handler {
	router := chi.NewRouter()
	router.Get("/health", healthHandler)

	return router
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
		log.Printf("ヘルスチェック応答の書き込みに失敗しました: %v", err)
	}
}
