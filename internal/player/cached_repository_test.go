package player

import (
	"context"
	"errors"
	"testing"
	"time"
)

type repositoryForCacheTest struct {
	player    Player
	findError error
	findCount int
}

func (repository *repositoryForCacheTest) Create(context.Context, string) (Player, error) {
	return Player{}, errors.New("このテストではCreateを呼び出しません")
}

func (repository *repositoryForCacheTest) FindByID(context.Context, string) (Player, error) {
	repository.findCount++
	return repository.player, repository.findError
}

type cacheForRepositoryTest struct {
	player      Player
	found       bool
	getError    error
	setPlayer   Player
	setDuration time.Duration
	setCount    int
	setError    error
}

func (cache *cacheForRepositoryTest) Get(context.Context, string) (Player, bool, error) {
	return cache.player, cache.found, cache.getError
}

func (cache *cacheForRepositoryTest) Set(_ context.Context, player Player, duration time.Duration) error {
	cache.setPlayer = player
	cache.setDuration = duration
	cache.setCount++
	return cache.setError
}

func TestCachedRepositoryReturnsCacheHitWithoutReadingPostgreSQL(t *testing.T) {
	cachedPlayer := Player{ID: "21faea44-59df-4ac7-b75b-c82d51a3de00", Name: "sumom"}
	source := &repositoryForCacheTest{}
	cache := &cacheForRepositoryTest{player: cachedPlayer, found: true}
	repository := NewCachedRepository(source, cache, 5*time.Minute)

	player, err := repository.FindByID(context.Background(), cachedPlayer.ID)
	if err != nil {
		t.Fatalf("Player取得に失敗しました: %v", err)
	}

	if player != cachedPlayer {
		t.Fatalf("Player = %#v, want %#v", player, cachedPlayer)
	}
	if source.findCount != 0 {
		t.Fatalf("PostgreSQLの取得回数 = %d, want 0", source.findCount)
	}
}

func TestCachedRepositoryReadsPostgreSQLAndWritesCacheOnMiss(t *testing.T) {
	foundPlayer := Player{ID: "21faea44-59df-4ac7-b75b-c82d51a3de00", Name: "sumom"}
	source := &repositoryForCacheTest{player: foundPlayer}
	cache := &cacheForRepositoryTest{}
	repository := NewCachedRepository(source, cache, 5*time.Minute)

	player, err := repository.FindByID(context.Background(), foundPlayer.ID)
	if err != nil {
		t.Fatalf("Player取得に失敗しました: %v", err)
	}

	if player != foundPlayer {
		t.Fatalf("Player = %#v, want %#v", player, foundPlayer)
	}
	if source.findCount != 1 {
		t.Fatalf("PostgreSQLの取得回数 = %d, want 1", source.findCount)
	}
	if cache.setCount != 1 {
		t.Fatalf("Redisの保存回数 = %d, want 1", cache.setCount)
	}
	if cache.setPlayer != foundPlayer {
		t.Fatalf("Redisへ保存するPlayer = %#v, want %#v", cache.setPlayer, foundPlayer)
	}
	if cache.setDuration != 5*time.Minute {
		t.Fatalf("Redisの有効期限 = %s, want %s", cache.setDuration, 5*time.Minute)
	}
}

func TestCachedRepositoryReadsPostgreSQLWhenRedisReturnsError(t *testing.T) {
	foundPlayer := Player{ID: "21faea44-59df-4ac7-b75b-c82d51a3de00", Name: "sumom"}
	source := &repositoryForCacheTest{player: foundPlayer}
	cache := &cacheForRepositoryTest{getError: errors.New("Redis connection error")}
	repository := NewCachedRepository(source, cache, 5*time.Minute)

	player, err := repository.FindByID(context.Background(), foundPlayer.ID)
	if err != nil {
		t.Fatalf("Player取得に失敗しました: %v", err)
	}

	if player != foundPlayer {
		t.Fatalf("Player = %#v, want %#v", player, foundPlayer)
	}
	if source.findCount != 1 {
		t.Fatalf("PostgreSQLの取得回数 = %d, want 1", source.findCount)
	}
	if cache.setCount != 0 {
		t.Fatalf("Redisの保存回数 = %d, want 0", cache.setCount)
	}
}
