package token

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
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
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	r.client.WithContext(ctx)
	return r.client.Set(key, value, expiration).Err()
}

func (r *RedisCache) Ping() error {
	return r.client.Ping().Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	r.client.WithContext(ctx)
	return r.client.Get(key).Result()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	r.client.WithContext(ctx)
	return r.client.Del(key).Err()
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

func (t *TokenRepository) SaveToken(ctx context.Context, token string, userID int) error {
	ttl, err := time.ParseDuration(t.cfg.TTL)
	err = t.redisClient.Set(ctx, fmt.Sprintf("user:%d", userID), token, ttl)
	return err
}

func (t *TokenRepository) GetUserIDByRefreshToken(ctx context.Context, token string) (string, error) {
	pattern := regexp.MustCompile(`^user:\d+$`)
	var cursor uint64
	for {
		keys, newCursor, err := t.redisClient.client.Scan(cursor, "user:*", 1).Result()
		if err != nil {
			log.Fatalf("Error scanning keys: %v", err)
		}

		for _, key := range keys {
			if pattern.MatchString(key) {
				value, err := t.redisClient.Get(ctx, key)
				if err != nil {
					log.Printf("Error getting key %s: %v", key, err)
					continue
				}

				if value == token {
					parts := strings.Split(key, ":")
					if len(parts) != 2 {
						return "", errors.New("invalid key format")
					}

					return parts[1], nil
				}
			}
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	return "", errors.New("token not found")
}

func (t *TokenRepository) CheckRefreshToken(ctx context.Context, key string) (string, error) {
	token, err := t.redisClient.Get(ctx, key)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token from cache")
	}

	return token, nil
}

func (t *TokenRepository) GetToken(ctx context.Context, key string) (string, error) {
	token, err := t.redisClient.Get(ctx, key)
	if err != nil {
		return "", errors.Wrap(err, "failed to get token from cache")
	}

	return token, nil
}

func (t *TokenRepository) DeleteToken(ctx context.Context, key string) error {
	err := t.redisClient.Del(ctx, key)
	if err != nil {
		return errors.Wrap(err, "failed to delete token from cache")
	}

	return nil
}
