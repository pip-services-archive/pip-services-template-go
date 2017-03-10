package api_v1

import (
    "net/http"
    "github.com/pip-services/pip-services-runtime-go"
    "github.com/pip-services/pip-services-runtime-go/api"            
    "github.com/pip-services/pip-services-template-go/db"
)

type DummyRestService struct {
    api.RestService
}

func NewDummyRestService(config *runtime.DynamicMap) *DummyRestService {
    defaultConfig := runtime.NewMapAndSet(
        "transport.type", "http",
        "transport.host", "localhost",
        //"transport.port", 3000,
        "transport.requestMaxSize", 1024 * 1024,
        "transport.connectTimeout", 60000,
        "transport.debug", true,
    )
    config = runtime.NewMapWithDefaults(config, defaultConfig)
        
    c := DummyRestService { RestService: *api.NewRestService("Dummy.RestService", config) }
    c.SetRegistrations(c.registrations)
    return &c
}

func (c *DummyRestService) Logic() db.IDummyDataAccess {
    return c.AbstractService.Logic().(db.IDummyDataAccess)
}

func (c *DummyRestService) registrations() {
    c.Register("GET", "/dummies", c.GetDummies)
    c.Register("GET", "/dummies/{dummy_id}", c.GetDummyById)
    c.Register("POST", "/dummies", c.CreateDummy)
    c.Register("PUT", "/dummies/{dummy_id}", c.UpdateDummy)
    c.Register("DELETE", "/dummies/{dummy_id}", c.DeleteDummy)
}

func (c *DummyRestService) GetDummies(w http.ResponseWriter, r *http.Request) {
    filter := c.GetFilterParams(r)
    paging := c.GetPagingParams(r)
    
    page, err := c.Logic().GetDummies(filter, paging)
    c.SendResult(w, page, err)        
}

func (c *DummyRestService) GetDummyById(w http.ResponseWriter, r *http.Request) {
    dummyID := c.GetRouteParam(r, "dummy_id")
    
    dummy, err := c.Logic().GetDummyById(dummyID)
    c.SendResult(w, dummy, err)        
}

func (c *DummyRestService) CreateDummy(w http.ResponseWriter, r *http.Request) {    
    dummyData := db.NewEmptyDummy()    
    if err := c.GetInputData(r, dummyData); err != nil {
        c.SendError(w, err)
        return
    }
    
    dummy, err := c.Logic().CreateDummy(dummyData)
    c.SendCreatedResult(w, dummy, err)        
}

func (c *DummyRestService) UpdateDummy(w http.ResponseWriter, r *http.Request) {
    dummyID := c.GetRouteParam(r, "dummy_id")
    dummyData := runtime.NewEmptyMap()
    if err := c.GetInputData(r, dummyData); err != nil {
        c.SendError(w, err)
        return
    }
    
    dummy, err := c.Logic().UpdateDummy(dummyID, dummyData)
    c.SendResult(w, dummy, err)        
}

func (c *DummyRestService) DeleteDummy(w http.ResponseWriter, r *http.Request) {
    dummyID := c.GetRouteParam(r, "dummy_id")

    err := c.Logic().DeleteDummy(dummyID)
    c.SendDeletedResult(w, err)        
}
