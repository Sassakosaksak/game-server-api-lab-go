package player

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{client: client}
}

func (cache *redisCache) Get(ctx context.Context, id string) (Player, bool, error) {
	value, err := cache.client.Get(ctx, playerCacheKey(id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return Player{}, false, nil
	}
	if err != nil {
		return Player{}, false, fmt.Errorf("RedisからPlayerを取得できませんでした: %w", err)
	}

	var player Player
	if err := json.Unmarshal(value, &player); err != nil {
		return Player{}, false, fmt.Errorf("RedisのPlayerデータを読み取れませんでした: %w", err)
	}

	return player, true, nil
}

func (cache *redisCache) Set(ctx context.Context, player Player, expiration time.Duration) error {
	value, err := json.Marshal(player)
	if err != nil {
		return fmt.Errorf("Redisへ保存するPlayerデータを作成できませんでした: %w", err)
	}

	if err := cache.client.Set(ctx, playerCacheKey(player.ID), value, expiration).Err(); err != nil {
		return fmt.Errorf("RedisへPlayerを保存できませんでした: %w", err)
	}

	return nil
}

func playerCacheKey(id string) string {
	return "player:" + id
}
