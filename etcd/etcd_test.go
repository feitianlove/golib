package etcd

import (
	"context"
	"fmt"
	"github.com/feitianlove/golib/config"
	"testing"
	"time"
)

func TestNewEtcdClient(t *testing.T) {
	conf := config.Config{
		Etcd: &config.Etcd{
			ListenPort:   "127.0.0.1:2379",
			TimeOut:      0,
			PrefixKey:    "",
			ProductKey:   "",
			BlackListKey: "",
		},
	}
	client, err := NewEtcdClient(&conf)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Close()
	}()
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(ctx, "/ftfeng/seckill/product", "robing")
	if err != nil {
		panic(err)
	}
	//cancel()
	resp, err := client.Get(ctx, "/ftfeng/seckill/product")
	if err != nil {
		panic(err)
	}
	//cancel()
	fmt.Println("------", resp)

}
