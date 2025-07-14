package data

import (
	"context"

	nacosregist "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
	demo1helloworld "github.com/orzkratos/demokratos/demo2kratos/api/helloworld/v1"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

type Demo1HttpClient struct {
	client        *http.Client
	greeterClient demo1helloworld.GreeterHTTPClient
}

func NewDemo1HttpClient(nacosNamingClient *NacosNamingClient, logger log.Logger) (*Demo1HttpClient, func()) {
	LOG := log.NewHelper(logger)

	// 这个写得非常好可以在更换时自动监听和更换IP地址，使用起来非常方便
	client := rese.P1(http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///demo1kratos.http"),
		http.WithDiscovery(nacosregist.New(nacosNamingClient.namingClient, nacosregist.WithGroup("demokratos"))),
		http.WithMiddleware(func(handler middleware.Handler) middleware.Handler {
			LOG.Infof("handle http request in middleware")
			return func(ctx context.Context, req any) (any, error) {
				// set auth info into context then request remote
				return handler(ctx, req)
			}
		}),
	))
	// 这里应该用 demo1helloworld 但实际上用的是 demo2helloworld 出于演示的目的只能这样吧
	// cp from https://github.com/go-kratos/kratos/blob/d6f5f00cf562b46322b0ed42d183b1b873c0a68f/transport/http/client_test.go#L339
	greeterClient := demo1helloworld.NewGreeterHTTPClient(client)
	cleanup := func() {
		must.Done(client.Close())
	}
	return &Demo1HttpClient{
		client:        client,
		greeterClient: greeterClient,
	}, cleanup
}

func (c *Demo1HttpClient) GetGreeterClient() demo1helloworld.GreeterHTTPClient {
	return c.greeterClient
}
