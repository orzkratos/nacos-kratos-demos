package data

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/yyle88/rese"
)

type NacosNamingClient struct {
	namingClient naming_client.INamingClient
}

func NewNacosNamingClient() *NacosNamingClient {
	// cp from https://github.com/go-kratos/examples/blob/61daed1ec4d5a94d689bc8fab9bc960c6af73ead/registry/nacos/client/main.go#L16
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/demo2kratos/log",
		CacheDir:            "/tmp/nacos/demo2kratos/cache",
		LogLevel:            "debug",
	}

	// cp from https://github.com/go-kratos/examples/blob/61daed1ec4d5a94d689bc8fab9bc960c6af73ead/registry/nacos/client/main.go#L31
	namingClient := rese.V1(clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	))
	return &NacosNamingClient{
		namingClient: namingClient,
	}
}
