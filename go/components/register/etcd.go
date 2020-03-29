package etcdv3

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"net"
	"strings"
)

var Deregister = make(chan struct{})

func Register(target, service, host, port string, ttl int64) error {
	// 拼接host:port
	serviceValue := net.JoinHostPort(host, port)
	serviceKey := service

	var err error

	cli, err := clientv3.New(clientv3.Config{
		Endpoints: strings.Split(target, ","),
	})

	if err != nil {
		return fmt.Errorf("create clientv3 failed:%v", err)
	}

	resp, err := cli.Grant(context.TODO(), ttl)

	if err != nil {
		return fmt.Errorf("lease failed:%v", err)
	}

	if _, err := cli.Put(context.TODO(), serviceKey, serviceValue, clientv3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf("set service %s failed:%v", service, err.Error())
	}

	if _, err := cli.KeepAlive(context.TODO(), resp.ID); err != nil {
		return fmt.Errorf("refresh service %s with ttl to clientv3 failed: %s", service, err.Error())
	}

	go func() {
		<-Deregister
		cli.Delete(context.Background(), serviceKey)
		Deregister <- struct{}{}
	}()

	return nil
}

func UnRegister() {
	Deregister <- struct{}{}
	<-Deregister
}
