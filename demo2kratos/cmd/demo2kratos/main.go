package main

import (
	"fmt"
	"os"

	nacosconfig "github.com/go-kratos/kratos/contrib/config/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/orzkratos/demokratos/demo2kratos/internal/conf"
	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern/zerotern"
	_ "go.uber.org/automaxprocs"
)

const nacosGroup = "demokratos"
const dataID = "demo2kratos.yaml" //需要带后缀，根据后缀选择解码器，假如不带后缀拿到的是Base64编码的数据
const defaultServiceName = "demo2kratos"

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
)

func init() {
	fmt.Println("service-name:", Name)
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(done.VCE(os.Hostname()).Omit()),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	// 有的时候会没有服务名称，需要默认值
	Name = zerotern.VV(Name, defaultServiceName)

	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", done.VCE(os.Hostname()).Omit(),
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	// cp from https://github.com/go-kratos/kratos/blob/d6f5f00cf562b46322b0ed42d183b1b873c0a68f/contrib/config/nacos/config_test.go#L16
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := &constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/demo2kratos/log",
		CacheDir:            "/tmp/nacos/demo2kratos/cache",
		LogLevel:            "debug",
	}

	configClient := rese.V1(clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	))

	// cp from https://github.com/go-kratos/kratos/blob/d6f5f00cf562b46322b0ed42d183b1b873c0a68f/contrib/config/nacos/config_test.go#L39
	source := nacosconfig.NewConfigSource(configClient, nacosconfig.WithGroup(nacosGroup), nacosconfig.WithDataID(dataID))

	c := config.New(
		config.WithSource(
			source,
		),
	)
	defer rese.F0(c.Close)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 假如需要随着配置更新而更新程序中的配置，就监听data字段的变动，因为data里基本是业务配置
	must.Done(c.Watch("data", func(s string, value config.Value) {
		must.Done(c.Scan(&bc))
		fmt.Println("config-data-update:", neatjsons.S(&bc))
	}))

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
