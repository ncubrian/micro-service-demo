package grpclb

import (
	"fmt"

	"github.com/CardInfoLink/log"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc/naming"
)

// watcher is the implementaion of grpc.naming.Watcher
type watcher struct {
	re            *resolver // re: Resolver
	conn          *zk.Conn
	addrs         []string
	isInitialized bool
}

// Close do nothing
func (w *watcher) Close() {
	w.conn.Close()
}

// Next to return the updates
func (w *watcher) Next() ([]*naming.Update, error) {
	// servicePath is the zookeeper dir to watch
	servicePath := fmt.Sprintf("/%s/%s", Prefix, w.re.serviceName)

	// query addresses from zookeeper
	addrs, _, childCh, err := w.conn.ChildrenW(servicePath)
	if err != nil {
		return nil, err
	}

	log.Infof("addrs: %v, chan: %p\n", addrs, childCh)

	// check if w is initialized
	if !w.isInitialized {
		w.isInitialized = true
		w.addrs = addrs

		//if not empty, return the updates or watcher new dir
		if l := len(addrs); l != 0 {
			updates := make([]*naming.Update, l)
			for i := range addrs {
				updates[i] = &naming.Update{Op: naming.Add, Addr: addrs[i]}
			}
			return updates, nil
		}
	}

	select {
	case childEvent := <-childCh:
		if childEvent.Type == zk.EventNodeDeleted {
			log.Info("receive znode delete event, ", childEvent)
		} else if childEvent.Type == zk.EventNodeChildrenChanged {
			log.Info("receive znode change event, ", childEvent)
			addrs, _, err := w.conn.Children(servicePath)
			if err != nil {
				return nil, err
			}
		
			log.Infof("addrs: %v\n", addrs)
			op, addr := w.processChildrenChanged(addrs)
			w.addrs = addrs
			return []*naming.Update{{Op: op, Addr: addr}}, nil
		} else {
			log.Infof("childEvent: %v\n", childEvent)
		}
	}

	return nil, nil
}

func (w *watcher) processChildrenChanged(addrs []string) (naming.Operation, string) {
	i := 0
	delIdx := 0
	for i = range addrs {
		j := 0
		for j < len(w.addrs) {
			if w.addrs[j] == addrs[i] {
				if j == delIdx {
					delIdx++
				}
				break
			}
			j++
		}

		if j == len(w.addrs) {
			return naming.Add, addrs[i]
		}
	}

	return naming.Delete, w.addrs[delIdx]
}
