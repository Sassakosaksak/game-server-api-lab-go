package player

import (
	"context"
	"time"
)

// CacheはPlayer取得結果を期限付きで保存・取得するための操作を定義する。
type Cache interface {
	Get(ctx context.Context, id string) (Player, bool, error)
	Set(ctx context.Context, player Player, expiration time.Duration) error
}
