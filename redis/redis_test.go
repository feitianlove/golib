package redis

import (
	"context"
	"fmt"
	"github.com/feitianlove/golib/config"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	conf := config.Config{
		Redis: &config.Redis{
			ListenPort:   ":6379",
			IdleTimeout:  0,
			MinIdleConns: 0,
			MaxConnAge:   0,
		},
	}
	client, _ := NewRedisClient(&conf)
	res, err := client.Do(context.Background(), "set", "ftfeng", "redis test").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
