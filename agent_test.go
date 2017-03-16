package agent

import (
	"fmt"
	"log"
	"testing"

	"github.com/smallnest/agent/pb"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

func createAgent() pb.AgentClient {
	//invoke service to agent
	conn, err := grpc.Dial("127.0.0.1:9981", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer conn.Close()

	return pb.NewAgentClient(conn)
}

func TestAgent_Call(t *testing.T) {
	beforeTest()

	agent := createAgent()

	//这个服务是在rpcx框架上实现的服务
	// 1. 服务的名称是MyService
	// 2. 要调用的方法是Handle
	// 3. 这个方法的输入参数类型是ProtoArgs
	// 4. 这个方法的输出参数类型是ProtoReply
	// 5. 数据通过protobuf序列化

	// 因此客户端可以使用gRPC支持的语言，比如java，python、C#， C++等调用rpcx的服务，只要:
	// 1. 输入输出都必须使用protobuf
	// 2. 客户端需要手工序列化/反序列化数据，而不是直接传入对象
	// 3. 需要指定要调用的服务和方法，格式为 "服务名.方法名"

	//准备要发送的数据，data是最终要调用的服务的输入参数，也就是MyService.Handle的输入参数，
	//已经序列化成[]byte
	data, _ := (&ProtoArgs{A: "world", B: 10}).Marshal()
	req := &pb.RpcRequest{Name: "MyService.Handle", Data: data}

	//通过代理调用服务
	reply, err := agent.Call(context.Background(), req)
	if err != nil {
		t.Errorf("failed to call rpcx service: %v", err)
	}

	//得到服务的返回结果reply.data，手工序列化成用户需要的结果ProtoReply
	data = reply.Data
	r := &ProtoReply{}
	err = r.Unmarshal(data)
	if err != nil {
		t.Errorf("failed to unmarhsal ProtoReply: %v", err)
	}

	// output result
	fmt.Printf("get result: %+v\n", r)

	afterTest()
}
