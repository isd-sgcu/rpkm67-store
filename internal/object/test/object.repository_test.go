package test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/isd-sgcu/rpkm67-store/config"
	"github.com/isd-sgcu/rpkm67-store/internal/object"
	httpClient "github.com/isd-sgcu/rpkm67-store/mocks/client/http"
	storeClient "github.com/isd-sgcu/rpkm67-store/mocks/client/store"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/suite"
)

type ObjectRepositoryTest struct {
	suite.Suite
	conf       *config.Store
	controller *gomock.Controller
	mockEndpoint string
}

func TestObjectRepository(t *testing.T) {
	suite.Run(t, new(ObjectRepositoryTest))
}

func (t *ObjectRepositoryTest) SetupTest() {
	t.conf = &config.Store{
		Endpoint: "mock-endpoint",
	}
	t.controller = gomock.NewController(t.T())
	t.mockEndpoint = "https://mock-endpoint/bucket/object"
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

func (t *ObjectRepositoryTest) TestDeleteError() {
	storeClient := storeClient.NewMockClient(t.controller)
	storeClient.EXPECT().RemoveObject(gomock.Any(), "bucket", "object", gomock.Any()).Return(errors.New("error"))

	repo := object.NewRepository(t.conf, storeClient, nil)

	err:= repo.Delete("bucket", "object")
	t.NotNil(err)
}

func (t *ObjectRepositoryTest) TestGetSuccess() {
	httpClient := httpClient.NewMockClient(t.controller)
	httpClient.EXPECT().Get(t.mockEndpoint).Return(&http.Response{
		StatusCode: http.StatusOK},nil)

	repo := object.NewRepository(t.conf, nil, httpClient)

	url,err:= repo.Get("bucket", "object")
	t.Nil(err)
	t.Assert().Equal(repo.GetURL("bucket","object"),url)
}

func (t *ObjectRepositoryTest) TestGetError() {
	httpClient := httpClient.NewMockClient(t.controller)
	httpClient.EXPECT().Get(t.mockEndpoint).Return(&http.Response{
		StatusCode: http.StatusOK},errors.New("error"))

	repo := object.NewRepository(t.conf, nil, httpClient)

	url,err:= repo.Get("bucket", "object")
	t.NotNil(err)
	t.Empty(url)
}

func (t *ObjectRepositoryTest) TestGetStatusNotOK() {
	httpClient := httpClient.NewMockClient(t.controller)
	httpClient.EXPECT().Get(t.mockEndpoint).Return(&http.Response{
		StatusCode: http.StatusNotFound},nil)

	repo := object.NewRepository(t.conf, nil, httpClient)

	url,err:= repo.Get("bucket", "object")
	t.Nil(err)
	t.Empty(url)
}

func (t *ObjectRepositoryTest) TestGetURL() {
	repo := object.NewRepository(t.conf, nil, nil)
	url := repo.GetURL("bucket","object")
	t.Assert().Equal(t.mockEndpoint,url)
}