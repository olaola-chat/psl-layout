package main

import (
	"flag"
	"os"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/olaola-chat/psl-be-partystar-pkg/i18n"
	"github.com/olaola-chat/psl-be-partystar-pkg/metrics"
	"github.com/olaola-chat/psl-be-partystar-pkg/nacos"
	pkgtracing "github.com/olaola-chat/psl-be-partystar-pkg/tracing"
	"github.com/olaola-chat/psl-be-partystar-pkg/util"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id = util.GenXID()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, reg *consul.Registry) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(reg),
	)
}

func main() {
	flag.Parse()
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	logger = log.NewFilter(logger, log.FilterLevel(log.ParseLevel(bc.GetLog().GetLevel())))
	log.SetLogger(logger)

	// init i18n config
	i18n.InitBundle("./configs/i18n")

	// init nacos client
	nacos.InitNacosClient(&nacos.NacosConfig{
		ServerIp:    bc.GetNacos().GetServerIp(),
		ServerPort:  bc.GetNacos().GetServerPort(),
		NamespaceId: bc.GetNacos().GetNamespaceId(),
		LogDir:      bc.GetNacos().GetLogDir(),
		CacheDir:    bc.GetNacos().GetCacheDir(),
		LogLevel:    bc.GetNacos().GetLogLevel(),
	})

	// init metrics
	metrics.InitGrpcMetrics(Name)

	// init tracing
	shutdown := pkgtracing.InitTracer(&pkgtracing.TracingConfig{
		ServiceName: bc.Server.Otel.ServiceName,
		Rate:        bc.Server.Otel.SampleRate,
		Endpoint:    bc.Server.Otel.Endpoint,
		Path:        bc.Server.Otel.Path,
	})
	defer shutdown()

	app, cleanup, err := wireApp(&bc, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
