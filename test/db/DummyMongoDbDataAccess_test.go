package db

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/pip-services/pip-services-runtime-go"    
    "github.com/pip-services/pip-services-template-go/db"    
)

type DummyMongoDbDataAccessTest struct {
    suite.Suite
    db *db.DummyMongoDbDataAccess
    fixture *DummyDataAccessFixture
}

func (suite *DummyMongoDbDataAccessTest) SetupSuite() {
    config := runtime.NewMapAndSet(
        "type", "mongodb",
        "uri", "mongodb://localhost/pipservicestest" )
    
    suite.db = db.NewDummyMongoDbDataAccess(config)
    suite.fixture = NewDummyDataAccessFixture(suite.db)
    
    err := suite.db.Init(runtime.NewReferences())
    assert.Nil(suite.T(), err)
    
    err = suite.db.Open()
    assert.Nil(suite.T(), err)
}

func (suite *DummyMongoDbDataAccessTest) SetupTest() {
    suite.db.Clear()
}

func (suite *DummyMongoDbDataAccessTest) TearDownSuite() {
    err := suite.db.Close()
    assert.Nil(suite.T(), err)
}

func (suite *DummyMongoDbDataAccessTest) TestCrudOperations() {
    suite.fixture.TestCrudOperations(suite.T())
}

func TestDummyMongoDbDataAccessTestSuite(t *testing.T) {
    suite.Run(t, new(DummyMongoDbDataAccessTest))
}