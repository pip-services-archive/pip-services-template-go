package db

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/pip-services/pip-services-runtime-go"    
    "github.com/pip-services/pip-services-template-go/db"    
)

type DummyFileDataAccessTest struct {
    suite.Suite
    db *db.DummyFileDataAccess
    fixture *DummyDataAccessFixture
}

func (suite *DummyFileDataAccessTest) SetupTest() {
    config := runtime.NewMapAndSet(
        "type", "file",
        "path", "../../dummies.json",
        "data", []*db.Dummy {} )
    
    suite.db = db.NewDummyFileDataAccess(config)
    suite.fixture = NewDummyDataAccessFixture(suite.db)
    
    err := suite.db.Init(runtime.NewReferences())
    assert.Nil(suite.T(), err)
    
    err = suite.db.Open()
    assert.Nil(suite.T(), err)
}

func (suite *DummyFileDataAccessTest) TearDownTest() {
    err := suite.db.Close()
    assert.Nil(suite.T(), err)
}

func (suite *DummyFileDataAccessTest) TestCrudOperations() {
    suite.fixture.TestCrudOperations(suite.T())
}

func (suite *DummyFileDataAccessTest) TestLoadData() {
    err := suite.db.Load()
    assert.Nil(suite.T(), err)
}

func TestDummyFileDataAccessTestSuite(t *testing.T) {
    suite.Run(t, new(DummyFileDataAccessTest))
}