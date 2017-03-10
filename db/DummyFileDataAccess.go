package db

import (
    "os"
    "bufio"
    "encoding/json"
    "github.com/pip-services/pip-services-runtime-go"
    "github.com/pip-services/pip-services-runtime-go/db"            
)

type DummyFileDataAccess struct {
    db.AbstractDataAccess
    
    path string
    maxPageSize int
    items []*Dummy
}

func NewDummyFileDataAccess(config *runtime.DynamicMap) *DummyFileDataAccess {
    defaultConfig := runtime.NewMapAndSet("maxPageSize", 100)
    config = runtime.NewMapWithDefaults(config, defaultConfig)
    
    c := DummyFileDataAccess { AbstractDataAccess: *db.NewAbstractDataAccess("Dummy.FileDataAccess", config) }

    return &c
}

func (c *DummyFileDataAccess) Init(refs *runtime.References) error {
    err := c.AbstractDataAccess.Init(refs)
    if err != nil { return err }
    
    if c.Config().HasNot("path") {
        return c.NewConfigError("Data file path is not set")
    }
    
    c.path = c.Config().GetString("path")
    c.maxPageSize = int(c.Config().GetInteger("maxPageSize"))
    
    return nil
}

func (c *DummyFileDataAccess) Open() error {
    c.items = nil
    
    // Initialize data from config
    data := c.Config().Get("data")
    if data != nil {
        items, ok := data.([]*Dummy)
        if ok == true { c.items = items }
    }
    
    if c.items == nil {
        if err := c.Load(); err != nil { return err }
    }
    
    return c.AbstractDataAccess.Open()
}

func (c *DummyFileDataAccess) Close() error {
    if err := c.Save(); err != nil { return err }
    return c.AbstractDataAccess.Close()
}

func (c *DummyFileDataAccess) Load() error {
    c.Trace("Loading data from file at " + c.path)
    
    // If doesn't exist then consider empty data
    if  _, err := os.Stat(c.path); os.IsNotExist(err) {
        c.items = make([]*Dummy, 0, 1000)
        return nil
    }
    
    f, err := os.Open(c.path)
    if err != nil {
        return c.NewReadError("Failed to read data file", err)
    }
    
    r := bufio.NewReader(f)
    decoder := json.NewDecoder(r)
    return decoder.Decode(&c.items)
}

func (c *DummyFileDataAccess) Save() error {
    c.Trace("Saving data to file at " + c.path)
    
    f, err := os.Create(c.path)
    if err != nil {
        return c.NewWriteError("Failed to write data file", err)
    }
    
    w := bufio.NewWriter(f)
    encoder := json.NewEncoder(w)    
    err = encoder.Encode(&c.items)
    if err != nil { return err }
    return w.Flush()
}

func (c *DummyFileDataAccess) Clear() error {
    c.items = make([]*Dummy, 0, 1000)
    return c.Save()
}

func (c *DummyFileDataAccess) GetDummies(filter *runtime.FilterParams, paging *runtime.PagingParams) (*DummyDataPage, error) {
    if filter == nil { filter = runtime.NewEmptyFilterParams() }
    
    var found []*Dummy = make([]*Dummy, 0, len(c.items))
    var key = (*runtime.DynamicMap)(filter).GetNullableString("key")
    
    for _, item := range c.items {
        if key != nil && *key != item.Key {
            continue
        }
            
        found = append(found, item)
    }
    
    if paging == nil { paging = runtime.NewEmptyPagingParams() }

    var total *int
    if paging.Paging { 
        t := len(found)
        total = &t    
    }
    
    skip := paging.GetSkip(-1)
    if skip > 0 {
        found = found[skip:]
    }

    take := paging.GetTake(c.maxPageSize)
    if int(take) < len(found) {
        found = found[0:take]
    }
    
    return NewDummyDataPage(total, found), nil
}

func (c *DummyFileDataAccess) GetDummyById(dummyID string) (*Dummy, error) {
    var found *Dummy
        
    for _, item := range c.items {
        if item.ID == dummyID {
            found = item
            break
        }
    }
     
    return found, nil
}

func (c *DummyFileDataAccess) CreateDummy(dummy *Dummy) (*Dummy, error) {
    if dummy.ID == "" { dummy.ID = c.CreateUUID() }

    c.items = append(c.items, dummy)

    return dummy, c.Save()
}

func (c *DummyFileDataAccess) UpdateDummy(dummyID string, dummy *runtime.DynamicMap) (*Dummy, error) {
    if dummy == nil { return nil, nil }
    
    var updated *Dummy    
    for _, item := range c.items {
        if item.ID == dummyID {
            updated = item
            dummy.AssignTo(updated)
            break
        }
    }
    
    if updated == nil { return nil, nil }    
    return updated, c.Save()
}

func (c *DummyFileDataAccess) DeleteDummy(dummyID string) error {
    items := make([]*Dummy, 0, len(c.items))
    
    var deleted *Dummy
    for _, item := range c.items {
        if item.ID != dummyID {
            items = append(items, item)
        } else {
            deleted = item
        }
    }
    
    if deleted == nil { return nil }        
    c.items = items    
    return c.Save()
}
