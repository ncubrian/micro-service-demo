package grpclb

import (
	"fmt"
	"strings"

	"github.com/CardInfoLink/log"
	etcd3 "github.com/coreos/etcd/clientv3"
	_ "github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"
)

// Prefix should start and end with no slash
var Prefix = "etcd3_naming"
var client *etcd3.Client
var serviceKey string
var leaseId etcd3.LeaseID

// Register regist configuration into etcd
func Register(name string, host string, port int, target string, ttl int) error {
	serviceValue := fmt.Sprintf("%s:%d", host, port)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceValue)
	log.Debugf("etcd serviceKey is %v, serviceValue is %v", serviceKey, serviceValue)
	
	// get endpoints for register dial address
	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(target, ","),
	})
	if err != nil {
		return fmt.Errorf("grpclb: create etcd3 client failed: %v", err)
	}

	// minimum lease TTL is ttl-second
	resp, _ := client.Grant(context.TODO(), int64(ttl))
	// should get first, if not exist, set it
	_, err = client.Get(context.Background(), serviceKey)
	// refresh set to true for not notifying the watcher
	if _, err = client.Put(context.Background(), serviceKey, serviceValue, etcd3.WithLease(resp.ID)); err != nil {
		log.Errorf("grpclb: refresh service '%s' with ttl to etcd3 failed: %s", name, err.Error())
		return err
	}
	_, err = client.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		log.Error(err)
		return err
	}
	leaseId = resp.ID

	return nil
}

// UnRegister delete registered service from etcd
func UnRegister() (err error) {
	if leaseId == 0 {
		return fmt.Errorf("invalid lease id")
	}

	// revoking lease expires the key attached to its lease ID
	_, err = client.Revoke(context.TODO(), leaseId)
	if err != nil {
		log.Error("grpclb: unregister '%s' failed: %s", serviceKey, err.Error())
	} else {
		log.Infof("grpclb: unregister '%s' ok.", serviceKey)
		leaseId = 0
	}
	return err
}
