package agent

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/smallnest/agent/codec"
	"github.com/smallnest/agent/pb"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/clientselector"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

var client *rpcx.Client
var server *grpc.Server

func StartAgent(addr string, registry string, opts []string) {
	client = createClientSelector(registry, opts...)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server = grpc.NewServer()
	pb.RegisterAgentServer(server, &agentServer{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Stop() {
	server.Stop()
}

type agentServer struct{}

func (s *agentServer) Call(ctx context.Context, req *pb.RpcRequest) (*pb.RpcReply, error) {
	//invoke services in rpcx
	var reply []byte
	err := client.Call(req.Name, req.Data, &reply)
	rpcReply := &pb.RpcReply{}
	if err == nil {
		rpcReply.Data = reply
	}
	return rpcReply, err
}

func createClientSelector(registry string, opts ...string) *rpcx.Client {
	var s rpcx.ClientSelector
	switch registry {
	case "zookeeper":
		s = clientselector.NewZooKeeperClientSelector(strings.Split(opts[0], ","), opts[1], 2*time.Minute, rpcx.WeightedRoundRobin, time.Minute)
	case "etcdv3":
		s = clientselector.NewEtcdV3ClientSelector(strings.Split(opts[0], ","), opts[1], 2*time.Minute, rpcx.WeightedRoundRobin, time.Minute)
	case "consul":
		s = clientselector.NewConsulClientSelector(opts[0], opts[1], 2*time.Minute, rpcx.WeightedRoundRobin, time.Minute)
	case "multi":
		servers := strings.Split(opts[0], ",")
		var serverPeers []*clientselector.ServerPeer
		for _, server := range servers {
			serverPeers = append(serverPeers, &clientselector.ServerPeer{Network: "tcp", Address: server})
		}
		s = clientselector.NewMultiClientSelector(serverPeers, rpcx.WeightedRoundRobin, time.Minute)
	case "direct":
		s = &rpcx.DirectClientSelector{Network: opts[0], Address: opts[1], DialTimeout: 10 * time.Second}
	}

	client := rpcx.NewClient(s)
	client.ClientCodecFunc = codec.NewClientCodec
	return client
}
