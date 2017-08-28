package main

import (
	"os"
	"strconv"
	"fmt"
	"net"
	"math/rand"
	"runtime"
	"strings"
	"time"
	"golang.org/x/net/context"

	"github.com/CardInfoLink/log"
	"rogers.chen/governor-helper-go/helper"
	"github.com/opentracing/opentracing-go/ext"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/ncubrian/micro-service-demo/service/pb"
)

type transServer struct{}
var randomSeed = rand.NewSource(time.Now().UnixNano())

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port, err := strconv.Atoi(os.Args[1])
	s, err := helper.Register("localhost", 23333, "servplat", "query", port)
	if err != nil {
		log.Errorf("failed to register, error is %v", err)
		return
	}

	pb.RegisterTransactionServer(s.GetServer(), &transServer{})

	s.Start()
}

func (s *transServer) Add(ctx context.Context, in *pb.Trans) (*pb.Resp, error) {
	log.Infof("Adding trans %#+v\n", in)
	return &pb.Resp{Ok: true}, nil
}

func (s *transServer) Update(ctx context.Context, in *pb.Trans) (*pb.Resp, error) {
	log.Infof("Updating trans %#+v\n", in)
	return &pb.Resp{Ok: true}, nil
}

func (s *transServer) Find(ctx context.Context, in *pb.QueryCond) (*pb.TransList, error) {
	log.Infof("Finding trans %#+v\n", in)
	list, err := dbMockFind(ctx)
	return list, err
}

func dbMockFind(ctx context.Context) (*pb.TransList, error) {
	// create new span using span found in context as parent (if none is found, our span becomes the trace root).
	resourceSpan, _ := opentracing.StartSpanFromContext(
		ctx,
		fmt.Sprintf("DB %s.dbMockFind", "trans"),
		// opentracing.StartTime(time.Now()),
	)
	defer func() {
		resourceSpan.Finish()
		// log.Debugf("%#+v", resourceSpan)
	}()
	// mark span as resource type
	ext.SpanKind.Set(resourceSpan, "resource")
	// name of the resource we try to reach
	ext.PeerService.Set(resourceSpan, "MongoDB")
	// hostname of the resource
	ext.PeerHostname.Set(resourceSpan, "transmgo1.showmoney.cn")
	// port of the resource
	ext.PeerPort.Set(resourceSpan, 27017)
	// let's binary annotate the query we run
	resourceSpan.SetTag(
		"query", "db.trans.find({})",
	)

	// Let's assume the query is going to take some time. Finding the right
	// world domination recipes is like searching for a needle in a haystack.
	time.Sleep(time.Duration(rand.New(randomSeed).Intn(100)) * time.Millisecond)

	return &pb.TransList{
		T: []*pb.Trans{
			&pb.Trans{Id: "id1", OrderNum: "order1", TransAmt: 1},
			&pb.Trans{Id: "id2", OrderNum: "order2", TransAmt: 2},
			&pb.Trans{Id: "id3", OrderNum: "order3", TransAmt: 3},
		},
	}, nil
}

// localIP 本机 IP
func localIPEx() (localIP string) {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		log.Error(err)
		localIP = "127.0.0.1"
	} else {
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
		conn.Close()
	}
	log.Infof("local ip:\t %s\n", localIP)
	return
}

// localIP 本机 IP
func localIP() (localIP string) {
	inter := "eth1"

	ifi, err := net.InterfaceByName(inter)
	if err != nil {
		log.Error(err)
		return localIPEx()
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		log.Error(err)
		return
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				break
			}
		}
	}
	// log.Debugf("local ip is %v", localIP)
	return
}
