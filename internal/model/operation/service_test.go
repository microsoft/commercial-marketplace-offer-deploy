package operation

import (
	"testing"

	"github.com/google/uuid"
	"github.com/microsoft/commercial-marketplace-offer-deploy/internal/model"
	"github.com/microsoft/commercial-marketplace-offer-deploy/test/fake"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type serviceTestSuite struct {
	suite.Suite
	db                 *fake.Database
	invokedOperationId uuid.UUID
}

func Test_OperationService(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

func (suite *serviceTestSuite) SetupSuite() {
	suite.db = fake.NewFakeDatabase(suite.T()).Setup(suite.setupDatabase)
}

func (suite *serviceTestSuite) Test_OperationService_deployment() {
	service := OperationService{
		db:  suite.db.Instance(),
		log: log.WithField("test", "Test_OperationService_deployment"),
	}
	operation, err := service.initialize(suite.invokedOperationId)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(operation)

	result := service.deployment()
	suite.Assert().NotNil(result)
}

func (suite *serviceTestSuite) setupDatabase(db *gorm.DB) {
	deployment := &model.Deployment{}
	db.Create(deployment)

	operation := &model.InvokedOperation{
		DeploymentId: deployment.ID,
		Name:         "test",
	}
	db.Create(operation)

	suite.invokedOperationId = operation.ID
}
