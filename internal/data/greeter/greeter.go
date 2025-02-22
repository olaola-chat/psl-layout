package greeter

import (
	"context"

	"psl-layout/internal/biz/greeter"

	"github.com/go-kratos/kratos/v2/log"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) greeter.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *greeterRepo) Save(ctx context.Context, g *greeter.Greeter) (*greeter.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) Update(ctx context.Context, g *greeter.Greeter) (*greeter.Greeter, error) {
	return g, nil
}

func (r *greeterRepo) FindByID(context.Context, int64) (*greeter.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListByHello(context.Context, string) ([]*greeter.Greeter, error) {
	return nil, nil
}

func (r *greeterRepo) ListAll(context.Context) ([]*greeter.Greeter, error) {
	return nil, nil
}
