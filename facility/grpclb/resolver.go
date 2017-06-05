package grpclb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc/naming"
)

// resolver is the implementaion of grpc.naming.Resolver
type resolver struct {
	serviceName string // service name to resolve
}

// NewResolver return resolver with service name
func NewResolver(serviceName string) *resolver {
	return &resolver{serviceName: serviceName}
}

// Resolve to resolve the service from etcd, target is the dial address of etcd
// target example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"
func (re *resolver) Resolve(target string) (naming.Watcher, error) {
	if re.serviceName == "" {
		return nil, errors.New("grpclb: no service name provided")
	}

	// generate zookeeper connection
	conn, connChan, err := zk.Connect(strings.Split(target, ","), 3 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("grpclb: creat zookeeper connection failed: %s", err.Error())
	}

	// 等待连接成功
	for {
		isConnected := false
		select {
		case connEvent := <-connChan:
			if connEvent.State == zk.StateConnected {
				isConnected = true
				fmt.Println("connect to zookeeper server success!")
			}
		case _ = <-time.After(time.Second * 3): // 3秒仍未连接成功则返回连接超时
			return nil, fmt.Errorf("connect to zookeeper server timeout!")
		}

		if isConnected {
			break
		}
	}

	// Return watcher
	return &watcher{re: re, conn: conn}, nil
}
