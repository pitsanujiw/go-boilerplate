package memstore

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"

	"github.com/pitsanujiw/go-boilerplate/config"
)

var ErrCacheMiss = errors.New("cache: key does not exist")

type CacheService interface {
	Set(ctx context.Context, key string, pValue any, ttl time.Duration) error

	Get(ctx context.Context, key string, pValue any) error

	Exists(ctx context.Context, key string) bool

	Delete(ctx context.Context, key string) error
}

type cacheService struct {
	rc    *redis.Client
	cache *cache.Cache
}

func New(cfg *config.App) CacheService {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &cacheService{
		rc: rc,
		cache: cache.New(&cache.Options{
			Redis:     rc,
			Marshal:   json.Marshal,
			Unmarshal: json.Unmarshal,
		}),
	}
}

func (c *cacheService) Set(ctx context.Context, key string, pValue interface{}, ttl time.Duration) error {
	err := c.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: pValue,
		TTL:   ttl,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *cacheService) Get(ctx context.Context, key string, pValue interface{}) error {
	err := c.cache.Get(ctx, key, pValue)
	switch {
	case errors.Is(err, cache.ErrCacheMiss):
		return ErrCacheMiss
	case err != nil:
		return err
	}
	return nil
}

func (c *cacheService) Exists(ctx context.Context, key string) bool {
	return c.cache.Exists(ctx, key)
}

func (c *cacheService) Delete(ctx context.Context, key string) error {
	if err := c.cache.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
