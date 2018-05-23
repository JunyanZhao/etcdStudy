package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("connect failed, err:", err)
		return
	}

	fmt.Println("connect succ")
	defer cli.Close()

	cli.Put(context.Background(), "/logagent/conf/", "8888888")
	go testEtcd(context.Background(), cli)
	go testEtcdGet(context.Background(), cli)
	for {
		rch := cli.Watch(context.Background(), "/logagent/conf/")
		rch = cli.Watch(context.Background(), "/wesure/conf/")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}
func testEtcd(ctx context.Context, client *clientv3.Client) {
	time.Sleep(time.Second * 6)
	client.Put(ctx, "/wesure/conf/", "999999")
}

func testEtcdGet(ctx context.Context, client *clientv3.Client) {
	time.Sleep(time.Second * 10)
	fmt.Println(client.Get(ctx, "/logagent/conf/"))
	fmt.Println(client.Get(ctx, "/wesure/conf/"))
	res, _ := client.Get(ctx, "/logagent/conf/")
	fmt.Println(string(res.Kvs[0].Value))
}
