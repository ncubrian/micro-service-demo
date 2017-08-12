package main

import (
	"os"
	"strconv"
	"fmt"
	"net"
	"runtime"
	"strings"
	"golang.org/x/net/context"

	"github.com/CardInfoLink/log"
	"rogers.chen/governor-helper/helper"

	pb "github.com/ncubrian/micro-service-demo/service/pb"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type transServer struct{}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port, err := strconv.Atoi(os.Args[1])
	err = helper.Register("localhost", 23333, "servplat", "query", port)
	if err != nil {
		log.Errorf("failed to register, error is %v", err)
		return
	}

	// Listen gRPC trans service port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Errorf("failed to listen, error is %v", err)
		return
	}


	// Init gRPC trans service
	s := grpc.NewServer()
	pb.RegisterTransactionServer(s, &transServer{})
	reflection.Register(s)
	if err = s.Serve(lis); err != nil {
		log.Errorf("failed to serve, error is %v", err)
		return
	}
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
