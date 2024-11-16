package cache

import (
	"context"
	"fmt"
	fb "github.com/Eugene-Usachev/fastbytes"
	loggerpkg "github.com/Eugune-Usachev/social-network/src/pkg/logger"
	"github.com/redis/rueidis"
	"time"
)

type RedisCache struct {
	client rueidis.Client
	logger loggerpkg.Logger
}

const (
	cacheDurationSeconds        = 300
	negativeCaseDurationSeconds = 300
)

func MustCreateRedisCache(addr, password string, logger loggerpkg.Logger) *RedisCache {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:   []string{addr},
		Password:      password,
		SelectDB:      0,
		MaxFlushDelay: 100 * time.Microsecond,
	})

	if err != nil {
		logger.Fatal(fmt.Sprintf("Error occurred when creating redis client: %s", err.Error()))
	}

	return &RedisCache{
		client: client,
		logger: logger,
	}
}

var _ Cache = (*RedisCache)(nil)

func (cache *RedisCache) IsNegativeCase(ctx context.Context, key string) bool {
	res, err := cache.client.Do(ctx, cache.client.B().Exists().Key(key).Build()).AsBool()
	if err != nil {
		cache.logger.Error(fmt.Sprintf("[Redis] Error occurred when checking negative case by key: %s, error: %s", key, err.Error()))
		return false
	}
	return res
}

func (cache *RedisCache) GetString(ctx context.Context, key string) (string, bool) {
	panic("TODO")
	//return cache.client.Do(ctx, cache.client.B().Get().Key(key).Build()).ToString()
}

func (cache *RedisCache) GetBytes(ctx context.Context, key string) ([]byte, bool) {
	panic("TODO")
	//return cache.client.Do(ctx, cache.client.B().Get().Key(key).Build()).AsBytes()
}

func (cache *RedisCache) SetString(ctx context.Context, key string, value string) {
	if err := cache.client.Do(ctx, cache.client.B().Set().Key(key).Value(value).ExSeconds(cacheDurationSeconds).Build()).Error(); err != nil {
		cache.logger.Error(fmt.Sprintf("[Redis] Error occurred when setting string by key: %s, error: %s", key, err.Error()))
	}
	return
}

func (cache *RedisCache) SetBytes(ctx context.Context, key string, value []byte) {
	if err := cache.client.Do(ctx, cache.client.B().Set().Key(key).Value(fb.B2S(value)).ExSeconds(cacheDurationSeconds).Build()).Error(); err != nil {
		cache.logger.Error(fmt.Sprintf("[Redis]  Erroroccurred when setting bytes by key: %s, error: %s", key, err.Error()))
	}
	return
}

func (cache *RedisCache) SetNegativeCase(ctx context.Context, key string) {
	if err := cache.client.Do(ctx, cache.client.B().Set().Key(key).Value("").ExSeconds(negativeCaseDurationSeconds).Build()).Error(); err != nil {
		cache.logger.Error(fmt.Sprintf("[Redis]  Erroroccurred when setting negative case by key: %s, error: %s", key, err.Error()))
	}
	return
}

func (cache *RedisCache) Delete(ctx context.Context, key string) error {
	return cache.client.Do(ctx, cache.client.B().Del().Key(key).Build()).Error()
}
