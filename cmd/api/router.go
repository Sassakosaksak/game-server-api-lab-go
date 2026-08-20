package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sassakosaksak/game-server-api-lab-go/internal/player"
	"github.com/go-chi/chi/v5"
)

type healthResponse struct {
	Status string `json:"status"`
}

// newRouterは本番起動とテストで同じルーティング設定を使うためにRouterを作る。
func newRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/health", healthHandler)

	return router
}

// registerPlayerRoutesはPlayer機能のHTTP処理をRouterへ登録する。
func registerPlayerRoutes(router chi.Router, handler *player.Handler) {
	router.Post("/players", handler.Create)
	router.Get("/players/{id}", handler.Get)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
		log.Printf("ヘルスチェック応答の書き込みに失敗しました: %v", err)
	}
}
