package server

import (
	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsHttpServer(addr string) *http.Server {
	s := http.NewServer(http.Address(addr))
	s.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	s.Handle("/metrics", promhttp.Handler())

	return s
}

func NewCouponGrpcMetricsHttpServer(bc *conf.Bootstrap) *http.Server {
	c := bc.GetServer()
	return NewMetricsHttpServer(c.GetGrpc().GetMetricsAddr())
}
