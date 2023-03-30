package data

import (
	"testing"
	"log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type uowSuite struct {
	suite.Suite	
	createEntity *Deployment
	updateEntity *Deployment
	deleteEntity *Deployment
	database Database
	uow UOW
}


func TestUnitOfWorkSuite(t *testing.T) {
	suite.Run(t, &uowSuite{})
}

func (s *uowSuite) SetupSuite() {
	c := &Deployment{Name: "CreateEntity"}
	s.createEntity = c

	u := &Deployment{Name: "UpdateEntity"}
	s.updateEntity = u

	d := &Deployment{Name: "DeleteEntity"}
	s.deleteEntity = d

	s.database = NewDatabase(&DatabaseOptions{UseInMemory: true})

	uow := NewUnitOfWork(s.database)
	s.uow = uow
}

func (s *uowSuite) SetupTest() {
	log.Print("SetupTest")
}

func (s *uowSuite) TestNewEntity(){
	s.uow.RegisterNew(s.createEntity)
	err := s.uow.Commit()
	require.NoError(s.T(), err)

	var result *Deployment
	s.database.Instance().Where("name = ?", s.createEntity.Name).First(&result)
	require.NotNil(s.T(), result)
	require.Equal(s.T(), s.createEntity.Name, result.Name)
}

func (s *uowSuite) TestUodateEntity(){

}

func (s *uowSuite) TestDeleteEntity(){
	s.uow.RegisterNew(s.deleteEntity)
	err := s.uow.Commit()
	require.NoError(s.T(), err)

	var insertResult *Deployment
	s.database.Instance().Where("name = ?", s.deleteEntity.Name).First(&insertResult)
	require.NotNil(s.T(), insertResult)
	require.Equal(s.T(), s.deleteEntity.Name, insertResult.Name)

	s.uow.RegisterDeleted(s.deleteEntity)
	err = s.uow.Commit()
	require.NoError(s.T(), err)

	var deleteResult *Deployment
	s.database.Instance().Where("name = ?", s.deleteEntity.Name).First(&deleteResult)
	require.NotNil(s.T(), insertResult)

	require.Equal(s.T(), uint(0), deleteResult.ID)
}