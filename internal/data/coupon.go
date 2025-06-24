package data

import (
	"layout/internal/biz"
	"layout/internal/data/coupon/gen"

	"github.com/go-kratos/kratos/v2/log"
)

type couponRepo struct {
	data *Data
	log  *log.Helper
	q    *gen.Query
}

func NewCouponRepo(data *Data, logger log.Logger) biz.CouponRepo {
	return &couponRepo{
		data: data,
		log:  log.NewHelper(logger),
		q:    gen.Use(data.XianshiDB),
	}
}
