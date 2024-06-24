package object

import (
	"context"

	proto "github.com/isd-sgcu/rpkm67-go-proto/rpkm67/store/object/v1"
	"go.uber.org/zap"
)

type Service interface {
	proto.ObjectServiceServer
}

type serviceImpl struct {
	proto.UnimplementedObjectServiceServer
	repo Repository
	log  *zap.Logger
	// client
}

func NewService(repo Repository, log *zap.Logger) proto.ObjectServiceServer {
	return &serviceImpl{repo: repo, log: log}
}

func (s *serviceImpl) Upload(_ context.Context, req *proto.UploadObjectRequest) (*proto.UploadObjectResponse, error) {
	return nil, nil
}

func (s *serviceImpl) FindByKey(_ context.Context, req *proto.FindByKeyObjectRequest) (*proto.FindByKeyObjectResponse, error) {
	return nil, nil
}

func (s *serviceImpl) DeleteByKey(_ context.Context, req *proto.DeleteByKeyObjectRequest) (*proto.DeleteByKeyObjectResponse, error) {
	return nil, nil
}
