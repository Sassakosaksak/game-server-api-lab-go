package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisOperationTimeoutは、キャッシュ障害でAPI応答を長く待たせないためのRedis操作ごとの上限時間。
const RedisOperationTimeout = 200 * time.Millisecond

// OpenRedisClientはRedis接続先の設定だけを行い、取得時の障害はPlayer取得処理でDBへフォールバックさせる。
func OpenRedisClient(address string) (*redis.Client, error) {
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("REDIS_ADDRが設定されていません")
	}

	return redis.NewClient(&redis.Options{
		Addr:                  address,
		DialTimeout:           RedisOperationTimeout,
		ReadTimeout:           RedisOperationTimeout,
		WriteTimeout:          RedisOperationTimeout,
		PoolTimeout:           RedisOperationTimeout,
		MaxRetries:            -1,
		DialerRetries:         1,
		ContextTimeoutEnabled: true,
	}), nil
}
