package test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/isd-sgcu/rpkm67-store/config"
	"github.com/isd-sgcu/rpkm67-store/internal/object"
	storeClient "github.com/isd-sgcu/rpkm67-store/mocks/client/store"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/suite"
)

type ObjectRepositoryTest struct {
	suite.Suite
	conf       *config.Store
	controller *gomock.Controller
}

func TestObjectRepository(t *testing.T) {
	suite.Run(t, new(ObjectRepositoryTest))
}

func (t *ObjectRepositoryTest) SetupTest() {
	t.conf = &config.Store{
		Endpoint: "mock-endpoint",
	}
	t.controller = gomock.NewController(t.T())
}

func (t *ObjectRepositoryTest) TestCreateObjectSuccess() {
	storeClient := storeClient.NewMockClient(t.controller)
	storeClient.EXPECT().
		PutObject(gomock.Any(), "mock-bucket", "mock-key", gomock.Any(), int64(0), gomock.Any()).
		Return(minio.UploadInfo{
			Key: "mock-key",
		}, nil)

	repo := object.NewRepository(t.conf, storeClient, nil)

	url, key, err := repo.Upload([]byte{}, "mock-bucket", "mock-key")
	t.Nil(err)
	t.Equal("mock-key", key)
	t.Equal(repo.GetURL("mock-bucket", "mock-key"), url)
}

func (t *ObjectRepositoryTest) TestUploadSuccess() {
	storeClient := storeClient.NewMockClient(t.controller)
	storeClient.EXPECT().PutObject(gomock.Any(), "bucket", "object", gomock.Any(), int64(0), gomock.Any()).Return(minio.UploadInfo{Key: "object"}, nil)

	repo := object.NewRepository(t.conf, storeClient, nil)

	url, key, err := repo.Upload([]byte{}, "bucket", "object")
	t.Nil(err)
	t.Equal("object", key)
	t.Equal(repo.GetURL("bucket", "object"), url)
}

func (t *ObjectRepositoryTest) TestUploadError() {
	storeClient := storeClient.NewMockClient(t.controller)
	storeClient.EXPECT().PutObject(gomock.Any(), "bucket", "object", gomock.Any(), int64(0), gomock.Any()).Return(minio.UploadInfo{}, errors.New("error"))

	repo := object.NewRepository(t.conf, storeClient, nil)

	url, key, err := repo.Upload([]byte{}, "bucket", "object")
	t.NotNil(err)
	t.Empty(url)
	t.Empty(key)
}

func (t *ObjectRepositoryTest) TestDeleteSuccess() {
	storeClient := storeClient.NewMockClient(t.controller)
	storeClient.EXPECT().RemoveObject(gomock.Any(), "bucket", "object", gomock.Any()).Return(nil)

	repo := object.NewRepository(t.conf, storeClient, nil)

	err:= repo.Delete("bucket", "object")
	t.Nil(err)
}

func (t *ObjectRepositoryTest) TestDeleteFail() {
	storeClient := storeClient.NewMockClient(t.controller)
	storeClient.EXPECT().RemoveObject(gomock.Any(), "bucket", "object", gomock.Any()).Return(errors.New("error"))

	repo := object.NewRepository(t.conf, storeClient, nil)

	err:= repo.Delete("bucket", "object")
	t.NotNil(err)
}







