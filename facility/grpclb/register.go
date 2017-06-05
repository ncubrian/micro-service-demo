package grpclb

import (
	"fmt"
	"strings"
	"time"

	"github.com/CardInfoLink/log"
	"github.com/samuel/go-zookeeper/zk"
)

// Prefix should start and end with no slash
var Prefix = "etcd3_naming"
var zkConn zk.Conn
var serviceKey string

// Register regist configuration into zookeeper
func Register(name string, host string, port int, target string, interval time.Duration) error {
	serviceValue := fmt.Sprintf("%s:%d", host, port)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceValue)
	
	// generate zookeeper connection
	conn, connChan, err := zk.Connect(strings.Split(target, ","), interval)
	if err != nil {
		return fmt.Errorf("grpclb: creat zookeeper connection failed: %s", err.Error())
	}

	// 等待连接成功
	for {
		isConnected := false
		select {
		case connEvent := <-connChan:
			if connEvent.State == zk.StateConnected {
				isConnected = true
				log.Info("connect to zookeeper server success!")
			}
		case _ = <-time.After(time.Second * 3): // 3秒仍未连接成功则返回连接超时
			return fmt.Errorf("connect to zookeeper server timeout!")
		}

		if isConnected {
			break
		}
	}

	rootDir := fmt.Sprintf("/%s", Prefix)
	// 判断zookeeper中是否存在root目录，不存在则创建该目录
	exist, _, err := conn.Exists(rootDir)
	if err != nil {
		return err
	}
	if !exist {
		path, err := conn.Create(rootDir, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
		if rootDir != path {
			return fmt.Errorf("Create returned different path, " + rootDir + " != " + path)
		}
	}

	serviceDir := fmt.Sprintf("%s/%s", rootDir, name)
	exist, _, err = conn.Exists(serviceDir)
	if err != nil {
		return err
	}
	if !exist {
		path, err := conn.Create(serviceDir, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
		if serviceDir != path {
			return fmt.Errorf("Create returned different path, " + serviceDir + " != " + path)
		}
	}

	exist, _, err = conn.Exists(serviceKey)
	if err != nil {
		return err
	}

	if !exist {
		path, err := conn.Create(serviceKey, ([]byte)(serviceValue), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
		if serviceKey != path {
			return fmt.Errorf("Create returned different path, " + serviceKey + " != " + path)
		}
	} else {
		return fmt.Errorf("service node %s already exist", serviceKey)
	}

	return nil
}
