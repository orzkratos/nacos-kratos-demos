package data

import (
	"context"
	"math/rand/v2"

	"github.com/go-kratos/kratos/v2/log"
	demo1helloworld "github.com/orzkratos/demokratos/demo2kratos/api/helloworld/v1"
	"github.com/orzkratos/demokratos/demo2kratos/internal/biz"
	"github.com/yyle88/erero"
)

type greeterRepo struct {
	data            *Data
	demo1GrpcClient *Demo1GrpcClient
	demo1HttpClient *Demo1HttpClient
	log             *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, demo1GrpcClient *Demo1GrpcClient, demo1HttpClient *Demo1HttpClient, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data:            data,
		demo1GrpcClient: demo1GrpcClient,
		demo1HttpClient: demo1HttpClient,
		log:             log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	// 这里是两个demo，grpc和http，随便用哪个都行，使用随机演示
	if rand.IntN(2) == 0 {
		// 这里应该用 demo1helloworld 但实际上用的是 demo2helloworld 出于演示的目的只能这样吧
		// cp from https://github.com/go-kratos/examples/blob/61daed1ec4d5a94d689bc8fab9bc960c6af73ead/registry/nacos/client/main.go#L52
		resp, err := r.demo1GrpcClient.greeterClient.SayHello(ctx, &demo1helloworld.HelloRequest{
			Name: g.Hello,
		})
		if err != nil {
			return nil, erero.Wro(err)
		}
		g.Hello = "message:[grpc-resp:" + resp.GetMessage() + "]"
	} else {
		// 这里应该用 demo1helloworld 但实际上用的是 demo2helloworld 出于演示的目的只能这样吧
		// cp from https://github.com/go-kratos/kratos/blob/d6f5f00cf562b46322b0ed42d183b1b873c0a68f/transport/http/client_test.go#L354
		resp, err := r.demo1HttpClient.greeterClient.SayHello(ctx, &demo1helloworld.HelloRequest{
			Name: g.Hello,
		})
		if err != nil {
			return nil, erero.Wro(err)
		}
		g.Hello = "message:[http-resp:" + resp.GetMessage() + "]"
	}
	return g, nil
}

func (r *greeterRepo) Update(ctx context.Context, g *biz.Greeter) (*biz.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) FindByID(context.Context, int64) (*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListByHello(context.Context, string) ([]*biz.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListAll(context.Context) ([]*biz.Greeter, error) {
	return nil, nil
}
