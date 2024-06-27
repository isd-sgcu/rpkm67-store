package object

import (
	"context"
	"fmt"

	proto "github.com/isd-sgcu/rpkm67-go-proto/rpkm67/store/object/v1"
	"github.com/isd-sgcu/rpkm67-store/config"
	"github.com/isd-sgcu/rpkm67-store/constant"
	"github.com/isd-sgcu/rpkm67-store/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	proto.ObjectServiceServer
}

type serviceImpl struct {
	proto.UnimplementedObjectServiceServer
	conf  *config.Store
	repo  Repository
	utils utils.Utils
	log   *zap.Logger
}

func NewService(repo Repository, conf *config.Store, log *zap.Logger, utils utils.Utils) Service {
	return &serviceImpl{
		repo:  repo,
		conf:  conf,
		utils: utils,
		log:   log,
	}
}

func (s *serviceImpl) Upload(_ context.Context, req *proto.UploadObjectRequest) (*proto.UploadObjectResponse, error) {
	randomString, err := s.utils.GenerateRandomString(10)
	if err != nil {
		s.log.Named("Upload").Error("GenerateRandomString: ", zap.Error(err))
		return nil, status.Error(codes.Internal, constant.InternalServerErrorMessage)
	}

	objectKey := req.Filename + "_" + randomString

	url, key, err := s.repo.Upload(req.Data, s.conf.BucketName, objectKey)
	if err != nil {
		s.log.Named("Upload").Error("Upload: ", zap.Error(err))
		return nil, status.Error(codes.Internal, constant.InternalServerErrorMessage)
	}

	return &proto.UploadObjectResponse{
		Object: &proto.Object{
			Url: url,
			Key: key,
		},
	}, nil
}

func (s *serviceImpl) FindByKey(_ context.Context, req *proto.FindByKeyObjectRequest) (*proto.FindByKeyObjectResponse, error) {
	if req.Key == "" {
		s.log.Named("FindByKey").Error("Key is empty")
		return nil, status.Error(codes.InvalidArgument, constant.KeyEmptyErrorMessage)
	}

	url, err := s.repo.Get(s.conf.BucketName, req.Key)
	if err != nil {
		s.log.Named("FindByKey").Error("Get: ", zap.Error(err))
		return nil, status.Error(codes.Internal, constant.InternalServerErrorMessage)
	}
	if url == "" {
		s.log.Named("FindByKey").Error(fmt.Sprintf("Object with key %v not found", req.Key))
		return nil, status.Error(codes.NotFound, constant.ObjectNotFoundErrorMessage)
	}

	return &proto.FindByKeyObjectResponse{
		Object: &proto.Object{
			Url: url,
			Key: req.Key,
		},
	}, nil
}

func (s *serviceImpl) DeleteByKey(_ context.Context, req *proto.DeleteByKeyObjectRequest) (*proto.DeleteByKeyObjectResponse, error) {
	if req.Key == "" {
		s.log.Named("DeleteByKey").Error("Key is empty")
		return &proto.DeleteByKeyObjectResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, constant.KeyEmptyErrorMessage)
	}

	err := s.repo.Delete(s.conf.BucketName, req.Key)
	if err != nil {
		s.log.Named("DeleteByKey").Error("Delete: ", zap.Error(err))
		return &proto.DeleteByKeyObjectResponse{
			Success: false,
		}, status.Error(codes.Internal, constant.InternalServerErrorMessage)
	}

	return &proto.DeleteByKeyObjectResponse{
		Success: true,
	}, nil
}

func (s *serviceImpl) GetURL(bucketName string, objectKey string) string {
	return "https://" + s.conf.Endpoint + "/" + bucketName + "/" + objectKey
}
