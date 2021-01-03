package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/feitianlove/golib/common/logger"
	"github.com/feitianlove/golib/config"
	"github.com/sirupsen/logrus"
	etcd "go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

type Etcd struct {
	ClientEtcd *etcd.Client
	RwLock     sync.RWMutex
}

func NewEtcdClient(conf *config.Config) (*etcd.Client, error) {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{conf.Etcd.ListenPort},
		DialTimeout: 5 * time.Second,
	})
	return cli, err
}
func (client *Etcd) WatchEtcdKey(key string, ctx context.Context) {
	c := client.ClientEtcd
	var secProductInfo []interface{}
	for {
		rch := c.Watch(ctx, key)
		var ReadSuccessFlag = true
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logger.Console.WithFields(logrus.Fields{
						"WatchEtcdKey": mvccpb.DELETE,
					}).Warn(fmt.Sprintf("Etcd delete key: %s", key))
				}
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						ReadSuccessFlag = false
						logger.Console.WithFields(logrus.Fields{
							"WatchEtcdKey": mvccpb.PUT,
						}).Error(fmt.Sprintf("json Unmarshal err: %s", err))
						continue
					}
				}
				logger.Console.WithFields(logrus.Fields{
					"WatchEtcdKey": ev.Type,
					"key":          ev.Kv.Key,
					"value":        ev.Kv.Value,
				}).Debug("Etcd key change")
			}
		}
		if ReadSuccessFlag {
			client.UpdateEtcdKey(secProductInfo)
		}
	}
}
func (client *Etcd) UpdateEtcdKey(data []interface{}) {
	var temp map[int]interface{} = make(map[int]interface{}, 10)
	for _, v := range data {
		//TODO 更新
		fmt.Println(v, temp)
	}
	client.RwLock.Lock()
	//更新
	client.RwLock.RLock()

}
