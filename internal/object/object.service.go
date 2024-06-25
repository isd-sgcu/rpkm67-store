package object

import (
	"bytes"
	"context"
	"net/http"

	proto "github.com/isd-sgcu/rpkm67-go-proto/rpkm67/store/object/v1"
	"github.com/isd-sgcu/rpkm67-store/config"
	"github.com/isd-sgcu/rpkm67-store/constant"
	"github.com/isd-sgcu/rpkm67-store/internal/client/store"
	"github.com/isd-sgcu/rpkm67-store/internal/model"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	proto.ObjectServiceServer
}

type serviceImpl struct {
	proto.UnimplementedObjectServiceServer
	conf        config.Config
	repo        Repository
	log         *zap.Logger
	storeClient store.Client
	httpClient  http.Client
}

func NewService(repo Repository, conf config.Config, storeClient store.Client, httpClient http.Client, log *zap.Logger) proto.ObjectServiceServer {
	return &serviceImpl{
		conf:        conf,
		repo:        repo,
		log:         log,
		storeClient: storeClient,
		httpClient:  httpClient,
	}
}

func (s *serviceImpl) Upload(_ context.Context, req *proto.UploadObjectRequest) (*proto.UploadObjectResponse, error) {
	if req.Data == nil {
		s.log.Named("Upload").Error(constant.FileNotFoundErrorMessage)
		return nil, status.Error(codes.NotFound, constant.FileNotFoundErrorMessage)
	}

	randomString, err := GenerateRandomString(10)
	if err != nil {
		s.log.Named("Upload").Error("GenerateRandomString: ", zap.Error(err))
	}

	objectKey := req.Filename + "_" + randomString
	buffer := bytes.NewReader(req.Data)

	uploadOutput, err := s.storeClient.Upload(s.conf.Store.BucketName, objectKey, buffer, buffer.Size(), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})

	if err != nil {
		s.log.Named("Upload").Error("Upload: ", zap.Error(err))
		return nil, status.Error(codes.Internal, constant.InternalServerErrorMessage)
	}

	objectResp := &proto.Object{
		Url: s.GetURL(s.conf.Store.BucketName, objectKey),
		Key: uploadOutput.Key,
	}

	err = s.repo.Upload(&model.Object{
		ImageUrl:  objectResp.Url,
		ObjectKey: objectResp.Key,
	})
	if err != nil {
		s.log.Named("Upload").Error("Upload: ", zap.Error(err))
		return nil, status.Error(codes.Internal, constant.InternalServerErrorMessage)
	}

	return &proto.UploadObjectResponse{
		Object: objectResp,
	}, nil
}

func (s *serviceImpl) FindByKey(_ context.Context, req *proto.FindByKeyObjectRequest) (*proto.FindByKeyObjectResponse, error) {
	return nil, nil
}

func (s *serviceImpl) DeleteByKey(_ context.Context, req *proto.DeleteByKeyObjectRequest) (*proto.DeleteByKeyObjectResponse, error) {
	if req.Key == "" {
		s.log.Named("DeleteByKey").Error("Key is empty")
		return &proto.DeleteByKeyObjectResponse{
			Success: false,
		}, status.Error(codes.InvalidArgument, constant.KeyEmptyErrorMessage)
	}

	err := s.storeClient.DeleteByKey(s.conf.Store.BucketName, req.Key, minio.RemoveObjectOptions{})
	if err != nil {
		s.log.Named("DeleteByKey").Error("DeleteByKey: ", zap.Error(err))
		return &proto.DeleteByKeyObjectResponse{
			Success: false,
		}, status.Error(codes.Internal, constant.InternalServerErrorMessage)
	}

	err = s.repo.DeleteByKey(req.Key)
	if err != nil {
		s.log.Named("DeleteByKey").Error("DeleteByKey: ", zap.Error(err))
	}

	return &proto.DeleteByKeyObjectResponse{
		Success: true,
	}, nil
}

func (s *serviceImpl) GetURL(bucketName string, objectKey string) string {
	return "https://" + s.conf.Store.Endpoint + "/" + bucketName + "/" + objectKey
}
