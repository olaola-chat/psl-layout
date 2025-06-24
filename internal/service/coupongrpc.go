package service

import (
	"context"
	"layout/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/olaola-chat/psl-be-protocol/coupon/v1"
)

type CouponGrpcService struct {
	pb.UnimplementedCouponGrpcServer

	uc  *biz.CouponUsecase
	log *log.Helper
}

func NewCouponGrpcService(logger log.Logger, uc *biz.CouponUsecase) *CouponGrpcService {
	return &CouponGrpcService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *CouponGrpcService) TryUse(ctx context.Context, req *pb.ReqCouponGrpcUse) (*pb.RespCouponGrpcUse, error) {
	return &pb.RespCouponGrpcUse{}, nil
}
func (s *CouponGrpcService) ConfirmUse(ctx context.Context, req *pb.ReqCouponGrpcUse) (*pb.RespCouponGrpcUse, error) {
	return &pb.RespCouponGrpcUse{}, nil
}
func (s *CouponGrpcService) CancelUse(ctx context.Context, req *pb.ReqCouponGrpcUse) (*pb.RespCouponGrpcUse, error) {
	return &pb.RespCouponGrpcUse{}, nil
}
