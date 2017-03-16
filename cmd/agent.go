package main

import (
	"flag"

	"strings"

	"github.com/smallnest/agent"
)

var (
	addr     = flag.String("addr", ":9981", "listen address")
	registry = flag.String("reg", "", "注册中心类型，支持direct,multi,zookeeper,etcdv3,consul等类型")
	opts     = flag.String("opts", "", "所需参数，不同的注册中心需要不同的参数，参数以空格分隔")
)

func main() {
	flag.Parse()

	agent.StartAgent(*addr, *registry, strings.Fields(*opts))
}
