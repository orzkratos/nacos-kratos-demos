package data

import (
	"context"

	nacosregist "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	demo1helloworld "github.com/orzkratos/demokratos/demo2kratos/api/helloworld/v1"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	grpc2 "google.golang.org/grpc"
)

type Demo1Client struct {
	conn          *grpc2.ClientConn
	greeterClient demo1helloworld.GreeterClient
}

func NewDemo1Client(nacosNamingClient *NacosNamingClient, logger log.Logger) (*Demo1Client, func()) {
	conn := rese.P1(grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///demo1kratos.grpc"),
		grpc.WithDiscovery(nacosregist.New(nacosNamingClient.namingClient, nacosregist.WithGroup("demokratos"))),
	))
	// 这里应该用 demo1helloworld 但实际上用的是 demo2helloworld 出于演示的目的只能这样吧
	// cp from https://github.com/go-kratos/examples/blob/61daed1ec4d5a94d689bc8fab9bc960c6af73ead/registry/nacos/client/main.go#L51
	greeterClient := demo1helloworld.NewGreeterClient(conn)
	cleanup := func() {
		must.Done(conn.Close())
	}
	return &Demo1Client{
		conn:          conn,
		greeterClient: greeterClient,
	}, cleanup
}
