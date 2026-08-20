package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sassakosaksak/game-server-api-lab-go/internal/cache"
	"github.com/Sassakosaksak/game-server-api-lab-go/internal/database"
	"github.com/Sassakosaksak/game-server-api-lab-go/internal/player"
)

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

	redisAddress := os.Getenv("REDIS_ADDR")
	redisClient, err := cache.OpenRedisClient(redisAddress)
	if err != nil {
		log.Fatalf("Redis Clientの準備に失敗しました: %v", err)
	}
	defer redisClient.Close()

	postgresRepository := player.NewPostgresRepository(pool)
	playerCache := player.NewRedisCache(redisClient)
	cachedRepository := player.NewCachedRepository(postgresRepository, playerCache, 5*time.Minute)

	router := newRouter()
	registerPlayerRoutes(router, player.NewHandler(cachedRepository))

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("APIサーバーをポート8080で起動します")
	log.Fatal(server.ListenAndServe())
}
