package token

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/qRe0/auth-api/configs"
)

type RedisCache struct {
	client *redis.Client
}

type CacheClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Ping() error
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	r.client.WithContext(ctx)
	return r.client.Set(key, value, expiration).Err()
}

func (r *RedisCache) Ping() error {
	return r.client.Ping().Err()
}

type TokenRepository struct {
	cfg         configs.JWTConfig
	redisClient RedisCache
}

func NewTokenRepo(cfg configs.JWTConfig, redisClient RedisCache) *TokenRepository {
	return &TokenRepository{
		cfg:         cfg,
		redisClient: redisClient,
	}
}

func Init(cfg configs.RedisConfig) (RedisCache, error) {
	const connectionStringTemplate = "%s:%s"
	connectionString := fmt.Sprintf(connectionStringTemplate, cfg.Host, cfg.Port)
	redisClient := redis.NewClient(&redis.Options{
		Addr: connectionString,
	})

	cacheWrapper := RedisCache{client: redisClient}
	err := cacheWrapper.Ping()
	if err != nil {
		return RedisCache{}, errors.Wrap(err, "failed to connect to redis cache")
	}

	return cacheWrapper, nil
}

func (t *TokenRepository) SaveToken(ctx context.Context, token string, userID string) error {
	ttl, err := time.ParseDuration(t.cfg.TTL)
	err = t.redisClient.Set(ctx, fmt.Sprintf("user:%s", userID), token, ttl)
	return err
}
