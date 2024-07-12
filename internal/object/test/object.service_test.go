package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/isd-sgcu/rpkm67-go-proto/rpkm67/store/object/v1"
	"github.com/isd-sgcu/rpkm67-store/constant"
	"github.com/isd-sgcu/rpkm67-store/internal/object"
	"github.com/isd-sgcu/rpkm67-store/internal/utils"
	mock_object "github.com/isd-sgcu/rpkm67-store/mocks/object"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/isd-sgcu/rpkm67-store/config"
)

type ObjectServiceTest struct {
	suite.Suite
	controller          *gomock.Controller
	conf                *config.Store
	logger              *zap.Logger
	uploadObjectRequest *proto.UploadObjectRequest
}

func TestObjectService(t *testing.T) {
	suite.Run(t, new(ObjectServiceTest))
}

func (t *ObjectServiceTest) SetupTest() {
	t.controller = gomock.NewController(t.T())
	t.logger = zap.NewNop()
	t.conf = &config.Store{
		BucketName: "mock-bucket",
		Endpoint:   "mock-endpoint",
	}
	t.uploadObjectRequest = &proto.UploadObjectRequest{
		Filename: "object",
		Data:     []byte("data"),
	}
}

func (t *ObjectServiceTest) TestUploadInternalError() {
	repo := mock_object.NewMockRepository(t.controller)
	repo.EXPECT().Upload(t.uploadObjectRequest.Data, t.conf.BucketName, gomock.Any()).Return("", "", fmt.Errorf("error"))

	svc := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	expectedErr := status.Error(codes.Internal, constant.InternalServerErrorMessage).Error()

	actual, err := svc.Upload(context.Background(), t.uploadObjectRequest)

	t.Nil(actual)
	t.EqualError(err, expectedErr)
}

func (t *ObjectServiceTest) TestUploadSuccess() {
	repo := mock_object.NewMockRepository(t.controller)
	repo.EXPECT().Upload(t.uploadObjectRequest.Data, t.conf.BucketName, gomock.Any()).Return("url", "key", nil)

	svc := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	expected := &proto.UploadObjectResponse{
		Object: &proto.Object{
			Key: "key",
			Url: "url",
		},
	}

	actual, err := svc.Upload(context.Background(), t.uploadObjectRequest)

	t.Nil(err)
	t.Equal(expected, actual)
}

func (t *ObjectServiceTest) TestFindByKeyEmptyError() {
	repo := mock_object.NewMockRepository(t.controller)
	srv := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	findByKeyInput := &proto.FindByKeyObjectRequest{
		Key: "",
	}

	expectedErr := status.Error(codes.InvalidArgument, constant.KeyEmptyErrorMessage).Error()

	actual, err := srv.FindByKey(context.Background(), findByKeyInput)

	t.Nil(actual)
	t.EqualError(err, expectedErr)
}

func (t *ObjectServiceTest) TestFindByKeyInternalError() {
	findByKeyInput := &proto.FindByKeyObjectRequest{
		Key: "key",
	}

	repo := mock_object.NewMockRepository(t.controller)
	repo.EXPECT().Get(t.conf.BucketName, findByKeyInput.Key).Return("", fmt.Errorf("error"))

	srv := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	expectedErr := status.Error(codes.Internal, constant.InternalServerErrorMessage).Error()

	actual, err := srv.FindByKey(context.Background(), findByKeyInput)

	t.Nil(actual)
	t.EqualError(err, expectedErr)
}

func (t *ObjectServiceTest) TestFindByKeyNotFoundError() {
	findByKeyInput := &proto.FindByKeyObjectRequest{
		Key: "key",
	}

	repo := mock_object.NewMockRepository(t.controller)
	repo.EXPECT().Get(t.conf.BucketName, findByKeyInput.Key).Return("", nil)

	srv := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	expectedErr := status.Error(codes.NotFound, constant.ObjectNotFoundErrorMessage).Error()

	actual, err := srv.FindByKey(context.Background(), findByKeyInput)

	t.Nil(actual)
	t.EqualError(err, expectedErr)
}

func (t *ObjectServiceTest) TestFindByKeySuccess() {
	findByKeyInput := &proto.FindByKeyObjectRequest{
		Key: "key",
	}

	repo := mock_object.NewMockRepository(t.controller)
	repo.EXPECT().Get(t.conf.BucketName, findByKeyInput.Key).Return("url", nil)

	srv := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	expected := &proto.FindByKeyObjectResponse{
		Object: &proto.Object{
			Key: findByKeyInput.Key,
			Url: "url",
		},
	}

	actual, err := srv.FindByKey(context.Background(), findByKeyInput)

	t.Nil(err)
	t.Equal(expected, actual)
}

func (t *ObjectServiceTest) TestDeleteByKeyEmptyError() {
	repo := mock_object.NewMockRepository(t.controller)
	srv := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	deleteByKeyInput := &proto.DeleteByKeyObjectRequest{
		Key: "",
	}

	expectedErr := status.Error(codes.InvalidArgument, constant.KeyEmptyErrorMessage).Error()

	actual, err := srv.DeleteByKey(context.Background(), deleteByKeyInput)

	t.Equal(actual.Success, false)
	t.EqualError(err, expectedErr)
}

func (t *ObjectServiceTest) TestDeleteByKeyInternalError() {
	deleteByKeyInput := &proto.DeleteByKeyObjectRequest{
		Key: "key",
	}

	repo := mock_object.NewMockRepository(t.controller)
	repo.EXPECT().Delete(t.conf.BucketName, deleteByKeyInput.Key).Return(fmt.Errorf("error"))

	srv := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	expectedErr := status.Error(codes.Internal, constant.InternalServerErrorMessage).Error()

	actual, err := srv.DeleteByKey(context.Background(), deleteByKeyInput)

	t.Equal(actual.Success, false)
	t.EqualError(err, expectedErr)
}

func (t *ObjectServiceTest) TestDeleteByKeySuccess() {
	deleteByKeyInput := &proto.DeleteByKeyObjectRequest{
		Key: "key",
	}

	repo := mock_object.NewMockRepository(t.controller)
	repo.EXPECT().Delete(t.conf.BucketName, deleteByKeyInput.Key).Return(nil)

	srv := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	actual, err := srv.DeleteByKey(context.Background(), deleteByKeyInput)

	t.Nil(err)
	t.Equal(actual.Success, true)
}
