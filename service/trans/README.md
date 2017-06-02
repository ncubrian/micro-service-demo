# Build & dispatch
Build trans service, and dispatch it to vagrant virtual boxex

    $ ./packageAndDispatch.sh

# Run
Run trans service on virtual boxex

    $ cd ../../vagrant

    $ vagrant ssh app1
    $ cd /opt/service
    $ nohup ./transService >> transService.log 2>&1 &
    $ exit

    $ vagrant ssh app2
    $ cd /opt/service
    $ nohup ./transService >> transService.log 2>&1 &
    $ exit

    $ vagrant ssh app3
    $ cd /opt/service
    $ nohup ./transService >> transService.log 2>&1 &
    $ exit

# service address in etcd
######key value
/etcd3_naming/trans/192.168.99.11:9001		192.168.99.11:9001
/etcd3_naming/trans/192.168.99.22:9001		192.168.99.22:9001
/etcd3_naming/trans/192.168.99.33:9001		192.168.99.33:9001
