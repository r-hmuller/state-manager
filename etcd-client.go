package main

import (
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func saveValue(registry Registry) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	cli, _ := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"127.0.0.1:2379"},
	})
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	key := registry.service + "-" + registry.uuid
	out, err := json.Marshal(registry)
	if err != nil {
		panic(err)
	}

	_, err = kv.Put(ctx, key, string(out))
	if err != nil {
		panic(err)
	}
}

func GetLatestCheckpoint(key string) string {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	cli, _ := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"127.0.0.1:2379"},
	})
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	opts := []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(0),
	}

	gr, _ := kv.Get(ctx, "key", opts...)
	serviceUuid := ""
	var latestDatetime int32
	for _, item := range gr.Kvs {
		var payload Registry
		json.Unmarshal(item.Value, &payload)
		if payload.datetime > latestDatetime {
			serviceUuid = payload.uuid
		}
	}

	return serviceUuid
}
