package registry

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/hashicorp/consul/api"
)

func NewConsulRegistry(bc *conf.Bootstrap) *consul.Registry {
	registryConf := bc.GetRegistry()
	client, err := api.NewClient(&api.Config{
		Address: registryConf.GetConsul().GetAddr(),
	})
	if err != nil {
		panic(err)
	}
	reg := consul.New(
		client,
		consul.WithHealthCheck(true),
		consul.WithHealthCheckInterval(10),
		consul.WithHeartbeat(true),
		consul.WithTimeout(time.Second*20),
	)
	return reg
}

func GetOneNode(ctx context.Context, reg *consul.Registry, serviceName string) (string, error) {
	ins, err := reg.GetService(ctx, serviceName)
	if err != nil {
		return "nil", err
	}
	if len(ins) == 0 {
		return "", errors.New("no service found")
	}
	// 使用随机选择的负载均衡策略
	index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(ins))
	node := ins[index]
	for _, v := range node.Endpoints {
		if strings.HasPrefix(v, "grpc://") {
			return strings.TrimPrefix(v, "grpc://"), nil
		}
	}
	return "", errors.New("no grpc endpoint found")
}
