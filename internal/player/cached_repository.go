package player

import (
	"context"
	"log"
	"time"
)

type cachedRepository struct {
	source     Repository
	cache      Cache
	expiration time.Duration
}

func NewCachedRepository(source Repository, cache Cache, expiration time.Duration) Repository {
	return &cachedRepository{
		source:     source,
		cache:      cache,
		expiration: expiration,
	}
}

func (repository *cachedRepository) Create(ctx context.Context, name string) (Player, error) {
	return repository.source.Create(ctx, name)
}

func (repository *cachedRepository) FindByID(ctx context.Context, id string) (Player, error) {
	cachedPlayer, found, err := repository.cache.Get(ctx, id)
	if err != nil {
		log.Printf("RedisからPlayerを取得できないためPostgreSQLを参照します: %v", err)
		return repository.source.FindByID(ctx, id)
	}

	if found {
		log.Printf("Redis HIT: player:%s", id)
		return cachedPlayer, nil
	}

	log.Printf("Redis MISS: player:%s", id)
	foundPlayer, err := repository.source.FindByID(ctx, id)
	if err != nil {
		return Player{}, err
	}

	if err := repository.cache.Set(ctx, foundPlayer, repository.expiration); err != nil {
		log.Printf("RedisへのPlayer保存に失敗しました: %v", err)
	}

	return foundPlayer, nil
}
