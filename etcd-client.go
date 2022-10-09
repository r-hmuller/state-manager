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
	kv := getEtcdClient()

	key := registry.Service + "-" + registry.Uuid
	out, err := json.Marshal(registry)
	if err != nil {
		panic(err)
	}

	_, err = kv.Put(ctx, key, string(out))
	if err != nil {
		panic(err)
	}
}

func SaveKV(service string, key string, value string) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	kv := getEtcdClient()

	_, err := kv.Put(ctx, service+"."+key, value)
	if err != nil {
		panic(err)
	}
}

func GetKV(service string, key string) string {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	kv := getEtcdClient()
	opts := getEtcdOpts()

	kvs, err := kv.Get(ctx, service+"."+key, opts...)
	if err != nil {
		panic(err)
	}
	for _, item := range kvs.Kvs {
		return string(item.Value)
	}
	return ""
}

func GetLatestCheckpoint(key string) string {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
	kv := getEtcdClient()
	opts := getEtcdOpts()

	gr, _ := kv.Get(ctx, key, opts...)
	serviceUuid := ""
	var latestDatetime int64
	for _, item := range gr.Kvs {
		var payload Registry
		json.Unmarshal(item.Value, &payload)
		if payload.Datetime > latestDatetime {
			serviceUuid = payload.Uuid
		}
	}

	return serviceUuid
}
func getEtcdClient() clientv3.KV {
	cli, _ := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"127.0.0.1:2379"},
	})
	defer cli.Close()
	return clientv3.NewKV(cli)
}

func getEtcdOpts() []clientv3.OpOption {
	return []clientv3.OpOption{
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(0),
	}
}
