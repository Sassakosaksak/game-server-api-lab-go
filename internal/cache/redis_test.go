package cache

import "testing"

func TestOpenRedisClientSetsShortTimeoutsWithoutRetries(t *testing.T) {
	client, err := OpenRedisClient("redis:6379")
	if err != nil {
		t.Fatalf("Redisクライアントの作成に失敗しました: %v", err)
	}
	defer client.Close()

	options := client.Options()
	if options.DialTimeout != RedisOperationTimeout {
		t.Fatalf("接続タイムアウト = %s, want %s", options.DialTimeout, RedisOperationTimeout)
	}
	if options.ReadTimeout != RedisOperationTimeout {
		t.Fatalf("読み取りタイムアウト = %s, want %s", options.ReadTimeout, RedisOperationTimeout)
	}
	if options.WriteTimeout != RedisOperationTimeout {
		t.Fatalf("書き込みタイムアウト = %s, want %s", options.WriteTimeout, RedisOperationTimeout)
	}
	if options.PoolTimeout != RedisOperationTimeout {
		t.Fatalf("接続プール待機時間 = %s, want %s", options.PoolTimeout, RedisOperationTimeout)
	}
	if options.MaxRetries != 0 {
		t.Fatalf("コマンド再試行回数 = %d, want 0", options.MaxRetries)
	}
	if options.DialerRetries != 1 {
		t.Fatalf("接続試行回数 = %d, want 1", options.DialerRetries)
	}
	if !options.ContextTimeoutEnabled {
		t.Fatal("Contextのタイムアウトが有効ではありません")
	}
}
