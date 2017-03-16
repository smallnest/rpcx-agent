package agent

import (
	"fmt"
	"time"

	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/codec"
)

var rpcxAddress = "127.0.0.1:8972"
var agentAddress = "127.0.0.1:9981"

var rpcxs *rpcx.Server

type MyService struct{}

func (s *MyService) Handle(args *ProtoArgs, reply *ProtoReply) error {
	fmt.Printf("args:%+v\n", args)
	reply.C = "hello " + args.A
	reply.D = int32(len(args.A)) * args.B
	return nil
}

func startRpcxServer() *rpcx.Server {
	server := rpcx.NewServer()
	server.ServerCodecFunc = codec.NewProtobufServerCodec
	server.RegisterName("MyService", new(MyService))

	server.Start("tcp", rpcxAddress)

	return server
}

func beforeTest() {
	// start RPCX server
	rpcxs = startRpcxServer()
	//start Agent
	go StartAgent(agentAddress, "direct", []string{"tcp", rpcxAddress})
	time.Sleep(1e9)
}

func afterTest() {
	//clean
	rpcxs.Close()
	Stop()
}
