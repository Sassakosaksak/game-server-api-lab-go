package player

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	cacheconfig "github.com/Sassakosaksak/game-server-api-lab-go/internal/cache"
)

func TestRedisCacheStopsWaitingWhenOperationTimeoutExpires(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:                  "redis:6379",
		MaxRetries:            -1,
		DialerRetries:         1,
		ContextTimeoutEnabled: true,
		// Redisへ接続できない状態を外部サービスなしで再現し、Contextの期限切れまで待機する。
		Dialer: func(ctx context.Context, _, _ string) (net.Conn, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		},
	})
	defer client.Close()

	startedAt := time.Now()
	_, _, err := NewRedisCache(client).Get(context.Background(), "player-id")
	elapsed := time.Since(startedAt)

	if err == nil {
		t.Fatal("Redis取得の失敗を検出できませんでした")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Redis取得エラー = %v, want context deadline exceeded", err)
	}
	maximumElapsed := cacheconfig.RedisOperationTimeout + 800*time.Millisecond
	if elapsed > maximumElapsed {
		t.Fatalf("Redis待機時間 = %s, want %s以内", elapsed, maximumElapsed)
	}
}
