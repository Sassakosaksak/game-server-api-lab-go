package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Sassakosaksak/game-server-api-lab-go/internal/database"
	"github.com/Sassakosaksak/game-server-api-lab-go/internal/player"
	"github.com/go-chi/chi/v5"
)

type healthResponse struct {
	Status string `json:"status"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URLが設定されていません")
	}

	pool, err := database.OpenPool(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("PostgreSQLへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	router := newRouter()
	registerPlayerRoutes(router, player.NewHandler(player.NewPostgresRepository(pool)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("APIサーバーをポート8080で起動します")
	log.Fatal(server.ListenAndServe())
}

// 本番起動とテストで同じルーティング設定を使うためにRouterを作る。
func newRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/health", healthHandler)

	return router
}

// APIの入口でルーティングをまとめ、Player機能のHTTP処理を登録する。
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
