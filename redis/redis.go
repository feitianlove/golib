package redis

import (
	"context"
	"github.com/feitianlove/golib/config"
	redis "github.com/go-redis/redis/v8"
	"time"
)

func NewRedisClient(conf *config.Config) (*redis.Client, error) {
	redisServer := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.ListenPort,
		Password:     "", // no password set
		DB:           0,  // use default DB
		MinIdleConns: conf.Redis.MinIdleConns,
		MaxConnAge:   time.Millisecond,
		IdleTimeout:  time.Microsecond,
	})
	_, err := redisServer.Do(context.Background(), "set", "ftfeng", "redis test").Result()
	if err != nil {
		return nil, err
	}
	return redisServer, nil
}
