package deps

import (    
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/pip-services/pip-services-runtime-go"    
    "github.com/pip-services/pip-services-runtime-go/build"    
    "github.com/pip-services/pip-services-template-go/db"    
    "github.com/pip-services/pip-services-template-go/logic"    
    "github.com/pip-services/pip-services-template-go/api/version1"    
    "github.com/pip-services/pip-services-template-go/deps/version1"    
)

type DummyRestClientTest struct {
    suite.Suite
    
    db *db.DummyFileDataAccess
    ctrl *logic.DummyController
    client *deps_v1.DummyRestClient
    api *api_v1.DummyRestService
    refs *runtime.References
    fixture *DummyClientFixture
}

func (suite *DummyRestClientTest) SetupSuite() {
    dbConfig := runtime.NewMapAndSet(
        "type", "file",
        "path", "../../../dummies.json",
        "data", []*db.Dummy {}, 
    )    
    suite.db = db.NewDummyFileDataAccess(dbConfig)    
    suite.ctrl = logic.NewDummyController(nil)


    apiConfig := runtime.NewMapAndSet(
        "type", "rest",
        "transport.type", "http",
        "transport.host", "localhost",
        "transport.port", 3000,
    )
    suite.client = deps_v1.NewDummyRestClient(apiConfig)
    suite.api = api_v1.NewDummyRestService(apiConfig)
    
    suite.fixture = NewDummyClientFixture(suite.client)
    
    suite.refs = runtime.NewReferences().WithDB(suite.db).WithCtrl(suite.ctrl).WithDep("dummies", suite.client).WithAPI(suite.api)
           
    err := build.LifeCycleManager.InitAndOpen(suite.refs)
    assert.Nil(suite.T(), err)
}

func (suite *DummyRestClientTest) SetupTest() {
    suite.db.Clear()
}

func (suite *DummyRestClientTest) TearDownSuite() {
    err := build.LifeCycleManager.Close(suite.refs)
    assert.Nil(suite.T(), err)
}

func (suite *DummyRestClientTest) TestCrudOperations() {
    suite.fixture.TestCrudOperations(suite.T())
}

func TestDummyRestClientTestSuite(t *testing.T) {
    suite.Run(t, new(DummyRestClientTest))
}