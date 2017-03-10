package deps_v1

import (
    "github.com/pip-services/pip-services-runtime-go"
    "github.com/pip-services/pip-services-runtime-go/deps"            
    "github.com/pip-services/pip-services-template-go/db"
)

type DummyRestClient struct {
    deps.RestClient
}

func NewDummyRestClient(config *runtime.DynamicMap) *DummyRestClient {
    return &DummyRestClient { RestClient: *deps.NewRestClient("Dummy.RestClient", config) }
}

func (c *DummyRestClient) GetDummies(filter *runtime.FilterParams, paging *runtime.PagingParams) (*db.DummyDataPage, error) {
    timing := c.Instrument("Dummy.GetDummies")
    defer func() { timing.EndTiming() }()
    
    outputPage := db.NewEmptyDummyDataPage()
    _, err := c.Call("GET", "/dummies" + c.GetPagingAndFilterParams(filter, paging), nil, outputPage)
    
    if err != nil { return nil, err }
           
    return outputPage, nil
}

func (c *DummyRestClient) GetDummyById(dummyID string) (*db.Dummy, error) {
    timing := c.Instrument("Dummy.GetDummyById")
    defer func() { timing.EndTiming() }()
    
    outputDummy := db.NewEmptyDummy()
    result, err := c.Call("GET", "/dummies/" + dummyID, nil, outputDummy)
    
    if err != nil { return nil, err }
    if result == nil { return nil, nil }
           
    return outputDummy, nil
}

func (c *DummyRestClient) CreateDummy(dummy *db.Dummy) (*db.Dummy, error) {
    timing := c.Instrument("Dummy.CreateDummy")
    defer func() { timing.EndTiming() }()

    outputDummy := db.NewEmptyDummy()
    _, err := c.Call("POST", "/dummies", dummy, outputDummy)
    
    if err != nil { return nil, err }
       
    return outputDummy, nil
}

func (c *DummyRestClient) UpdateDummy(dummyID string, dummy *runtime.DynamicMap) (*db.Dummy, error) {
    timing := c.Instrument("Dummy.UpdateDummy")
    defer func() { timing.EndTiming() }()
    
    outputDummy := db.NewEmptyDummy()
    result, err := c.Call("PUT", "/dummies/" + dummyID, dummy, outputDummy)
    
    if err != nil { return nil, err }
    if result == nil { return nil, nil }
           
    return outputDummy, nil
}

func (c *DummyRestClient) DeleteDummy(dummyID string) error {
    timing := c.Instrument("Dummy.DeleteDummy")
    defer func() { timing.EndTiming() }()
    
    _, err := c.Call("DELETE", "/dummies/" + dummyID, nil, nil)
               
    return err
}
