package caching

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rahul-aut-ind/service-user/internal/config"
	"github.com/rahul-aut-ind/service-user/pkg/logger"
)

type (
	CacheHandler interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value string, ttl time.Duration) error
		Delete(ctx context.Context, key string) error
	}

	RedisClient struct {
		redisClient *redis.Client
		log         *logger.Logger
	}
)

const (
	KeyDoesNotExist = "key doesnot exist"
	DefaultTTL      = 30 * time.Second
)

func New(env *config.Env, l *logger.Logger) *RedisClient {
	ch := &RedisClient{
		log: l,
	}
	ch.redisClient = ch.initRedis(env.RedisAddress)
	return ch
}

func (rc *RedisClient) initRedis(addr string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if _, err := redisClient.Ping(context.TODO()).Result(); err != nil {
		rc.log.Fatalf("could not connect to redis: %v", err)
	}
	rc.log.Debug("connected to redis...")
	return redisClient
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	res, err := rc.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		rc.log.Debugf("%s %s", key, KeyDoesNotExist)
		return "", fmt.Errorf("%s", KeyDoesNotExist)
	}
	rc.log.Debugf("success reading from redis...")
	return res, nil
}

func (rc *RedisClient) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	_, err := rc.redisClient.Set(ctx, key, value, ttl).Result()
	if err != nil {
		rc.log.Infof("error setting %s :: err %s", key, err)
		return fmt.Errorf("err:: %s", err)
	}
	rc.log.Debugf("success added to redis...")
	return nil
}

func (rc *RedisClient) Delete(ctx context.Context, key string) error {
	_, err := rc.redisClient.Del(ctx, key).Result()
	if err != nil {
		rc.log.Infof("error deleting %s :: err %s", key, err)
		return fmt.Errorf("err:: %s", err)
	}
	rc.log.Debugf("success deleted from redis...")
	return nil
}
