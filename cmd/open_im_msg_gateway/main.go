package main

import (
	"flag"
	"fmt"
	"micro_servers/internal/msg_gateway/gate"
	"micro_servers/pkg/common/config"
	"micro_servers/pkg/common/constant"
	"micro_servers/pkg/common/log"
	"sync"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	defaultRpcPorts := config.Config.RpcPort.OpenImMessageGatewayPort
	defaultWsPorts := config.Config.LongConnSvr.WebsocketPort
	rpcPort := flag.Int("rpc_port", defaultRpcPorts[0], "rpc listening port")
	wsPort := flag.Int("ws_port", defaultWsPorts[0], "ws listening port")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("start rpc/msg_gateway server, port: ", *rpcPort, *wsPort)
	gate.Init(*rpcPort, *wsPort)
	gate.Run()
	wg.Wait()
}
