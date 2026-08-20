package cache

import (
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

// OpenRedisClientはRedis接続先の設定だけを行い、取得時の障害はPlayer取得処理でDBへフォールバックさせる。
func OpenRedisClient(address string) (*redis.Client, error) {
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("REDIS_ADDRが設定されていません")
	}

	return redis.NewClient(&redis.Options{Addr: address}), nil
}
