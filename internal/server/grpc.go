package server

import (
	"layout/internal/conf"
	"layout/internal/service"

	validate "github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	pkgMetrics "github.com/olaola-chat/psl-be-partystar-pkg/metrics"
	"github.com/olaola-chat/psl-be-partystar-pkg/middleware/metrics"
	coupon "github.com/olaola-chat/psl-be-protocol/coupon/v1"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(bc *conf.Bootstrap, couponSvc *service.CouponGrpcService, logger log.Logger) *grpc.Server {
	c := bc.Server
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			metrics.GrpcServer(
				"partystar-server",
				c.GetEnv(),
				metrics.WithRequests(pkgMetrics.GrpcServerRequestsCounter),
				metrics.WithMillSeconds(pkgMetrics.GrpcServerRequestMillSecondsHistogram),
			),
			logging.Server(logger),
			metadata.Server(),
			validate.ProtoValidate(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	coupon.RegisterCouponGrpcServer(srv, couponSvc)
	return srv
}
