package test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/isd-sgcu/rpkm67-store/config"
)

type ObjectServiceTest struct {
	suite.Suite
	controller *gomock.Controller
	conf       *config.Store
	logger     *zap.Logger
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
}
