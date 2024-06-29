package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	proto "github.com/isd-sgcu/rpkm67-go-proto/rpkm67/store/object/v1"
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
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	repo.EXPECT().Upload(t.uploadObjectRequest.Data, t.conf.BucketName, gomock.Any()).Return("", "", fmt.Errorf("error"))

	_, err := service.Upload(context.Background(), t.uploadObjectRequest)

	t.EqualError(err, "rpc error: code = Internal desc = Internal server error")
}

func (t *ObjectServiceTest) TestUpload() {
	repo := mock_object.NewMockRepository(t.controller)
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	repo.EXPECT().Upload(t.uploadObjectRequest.Data, t.conf.BucketName, gomock.Any()).Return("url", "key", nil)

	_, err := service.Upload(context.Background(), t.uploadObjectRequest)

	if err != nil {
		fmt.Println(err)
	}
}

func (t *ObjectServiceTest) TestFindByKeyEmptyError() {
	repo := mock_object.NewMockRepository(t.controller)
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	findByKeyInput := &proto.FindByKeyObjectRequest{
		Key: "",
	}

	_, err := service.FindByKey(context.Background(), findByKeyInput)

	t.EqualError(err, "rpc error: code = InvalidArgument desc = Key is empty")
}

func (t *ObjectServiceTest) TestFindByKeyInternalError() {
	repo := mock_object.NewMockRepository(t.controller)
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	findByKeyInput := &proto.FindByKeyObjectRequest{
		Key: "key",
	}

	repo.EXPECT().Get(t.conf.BucketName, findByKeyInput.Key).Return("", fmt.Errorf("error"))

	_, err := service.FindByKey(context.Background(), findByKeyInput)

	t.EqualError(err, "rpc error: code = Internal desc = Internal server error")
}

func (t *ObjectServiceTest) TestFindByKeyNotFoundError() {
	repo := mock_object.NewMockRepository(t.controller)
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	findByKeyInput := &proto.FindByKeyObjectRequest{
		Key: "key",
	}

	repo.EXPECT().Get(t.conf.BucketName, findByKeyInput.Key).Return("", nil)

	_, err := service.FindByKey(context.Background(), findByKeyInput)

	t.EqualError(err, "rpc error: code = NotFound desc = Object not found")
}

func (t *ObjectServiceTest) TestFindByKey() {
	repo := mock_object.NewMockRepository(t.controller)
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	findByKeyInput := &proto.FindByKeyObjectRequest{
		Key: "key",
	}

	repo.EXPECT().Get(t.conf.BucketName, findByKeyInput.Key).Return("url", nil)

	_, err := service.FindByKey(context.Background(), findByKeyInput)

	if err != nil {
		fmt.Println(err)
	}
}

func (t *ObjectServiceTest) TestDeleteByKeyEmptyError() {
	repo := mock_object.NewMockRepository(t.controller)
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	deleteByKeyInput := &proto.DeleteByKeyObjectRequest{
		Key: "",
	}

	_, err := service.DeleteByKey(context.Background(), deleteByKeyInput)

	t.EqualError(err, "rpc error: code = InvalidArgument desc = Key is empty")
}

func (t *ObjectServiceTest) TestDeleteByKeyInternalError() {
	repo := mock_object.NewMockRepository(t.controller)
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	deleteByKeyInput := &proto.DeleteByKeyObjectRequest{
		Key: "key",
	}

	repo.EXPECT().Delete(t.conf.BucketName, deleteByKeyInput.Key).Return(fmt.Errorf("error"))

	_, err := service.DeleteByKey(context.Background(), deleteByKeyInput)

	t.EqualError(err, "rpc error: code = Internal desc = Internal server error")
}

func (t *ObjectServiceTest) TestDeleteByKey() {
	repo := mock_object.NewMockRepository(t.controller)
	service := object.NewService(repo, t.conf, t.logger, utils.NewRandomUtils())

	deleteByKeyInput := &proto.DeleteByKeyObjectRequest{
		Key: "key",
	}

	repo.EXPECT().Delete(t.conf.BucketName, deleteByKeyInput.Key).Return(nil)

	_, err := service.DeleteByKey(context.Background(), deleteByKeyInput)

	if err != nil {
		fmt.Println(err)
	}
}
