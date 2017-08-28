package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"
	_ "time"

	"github.com/CardInfoLink/log"
	"rogers.chen/governor-helper-go/helper"
	pb "github.com/ncubrian/micro-service-demo/service/pb"
)

var (
	transServ = flag.String("trans service", "trans", "transaction service name")
	zkReg     = flag.String("reg", "192.168.1.213:2181,192.168.1.224:2181,192.168.1.226:2181", "register zookeeper address")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	// Init grpc load balancer
	conn, cancel, pctx, err := helper.Discover("localhost", 23333, "servplat", "query", "consumer", 32222)
	if err != nil {
		log.Error(err)
		return
	}
	ctx := *pctx

	defer cancel()
	defer conn.Close()
	
	transClient := pb.NewTransactionClient(conn)
	bio := bufio.NewReader(os.Stdin)
	
	for {
		line, _, err := bio.ReadLine()
		if string(line) == "exit" {
			break
		}

		var retStr string
		transList, err := transClient.Find(ctx, &pb.QueryCond{})
		if err != nil {
			retStr = fmt.Sprintf("failed to call find on trans server, error is %v", err)
			log.Errorf(retStr)
			return
		}
	
		transSlice := transList.GetT()
		if transSlice == nil {
			retStr = fmt.Sprintf("nil slice returned from trans server")
			log.Errorf(retStr)
			return
		}
	
		for _, t := range transSlice {
			retStr += fmt.Sprintf("\nid: %s, orderNum: %s, transAmt: %d", t.GetId(), t.GetOrderNum(), t.GetTransAmt())
		}
		log.Infof(retStr)
	}
}
