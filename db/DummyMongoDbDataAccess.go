package db

import (
	"github.com/pip-services/pip-services-runtime-go"
	"github.com/pip-services/pip-services-runtime-go/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type DummyMongoDbDataAccess struct {
	db.AbstractDataAccess

	url          string
	databaseName string
	maxPageSize  int
	session      *mgo.Session
	database     *mgo.Database
	collection   *mgo.Collection
}

func NewDummyMongoDbDataAccess(config *runtime.DynamicMap) *DummyMongoDbDataAccess {
	defaultConfig := runtime.NewMapAndSet(
		//"uri": nil,
		"options.server.poolSize", 4,
		"options.server.socketOptions.keepAlive", 1,
		"options.server.socketOptions.connectTimeoutMS", 5000,
		"options.server.auto_connect", true,
		"maxPageSize", 100,
		"debug", true)
	config = runtime.NewMapWithDefaults(config, defaultConfig)

	c := DummyMongoDbDataAccess{AbstractDataAccess: *db.NewAbstractDataAccess("Dummy.FileDataAccess", config)}

	return &c
}

func (c *DummyMongoDbDataAccess) Init(refs *runtime.References) error {
	err := c.AbstractDataAccess.Init(refs)
	if err != nil {
		return err
	}

	if c.Config().HasNot("uri") {
		return c.NewConfigError("MongoDB uri is not set")
	}

	c.maxPageSize = int(c.Config().GetInteger("maxPageSize"))

	return nil
}

func (c *DummyMongoDbDataAccess) Open() error {
	uri := c.Config().GetString("uri")
	lastIndex := strings.LastIndex(uri, "/")
	c.url = uri[0:lastIndex]
	c.databaseName = uri[lastIndex+1:]

	if c.url == "" {
		c.NewConfigError("Wrong mongodb uri").WithDetails(uri)
	}
	if c.databaseName == "" {
		c.NewConfigError("Database name in the uri is not set").WithDetails(uri)
	}

	c.Trace("Connecting to mongodb at " + uri)
	var err error
	c.session, err = mgo.Dial(uri)
	if err != nil {
		return c.NewOpenError("Failed to connect to MongoDB", err)
	}

	c.database = c.session.DB(c.databaseName)
	c.collection = c.database.C("dummies")

	return c.AbstractDataAccess.Open()
}

func (c *DummyMongoDbDataAccess) Close() error {
	if c.session != nil {
		c.session.Close()
		c.session = nil
		c.database = nil
		c.collection = nil
	}

	return c.AbstractDataAccess.Close()
}

func (c *DummyMongoDbDataAccess) Clear() {
	if c.database == nil {
		panic("Not connected to MongoDB")
	}

	if err := c.database.DropDatabase(); err != nil {
		panic(err.Error())
	}
}

func (c *DummyMongoDbDataAccess) GetDummies(filter *runtime.FilterParams, paging *runtime.PagingParams) (*DummyDataPage, error) {
	if filter == nil {
		filter = runtime.NewEmptyFilterParams()
	}
	if paging == nil {
		paging = runtime.NewEmptyPagingParams()
	}

	f := (*runtime.DynamicMap)(filter)
	condition := bson.M{}
	if f.Has("key") {
		condition["key"] = f.GetString("key")
	}

	query := c.collection.Find(condition)

	var total *int
	if paging.Paging {
		t, err := query.Count()
		if err != nil {
			c.NewReadError("Failed to read query count", err)
		}
		total = &t
	}

	skip := paging.GetSkip(-1)
	if skip > 0 {
		query = query.Skip(skip)
	}

	take := paging.GetTake(c.maxPageSize)
	query = query.Limit(take)

	found := make([]*Dummy, 0, c.maxPageSize)
	if err := query.All(&found); err != nil {
		c.NewReadError("Failed to read query", err)
	}

	return NewDummyDataPage(total, found), nil
}

func (c *DummyMongoDbDataAccess) GetDummyById(dummyID string) (*Dummy, error) {
	found := make([]*Dummy, 0, 1)

	err := c.collection.Find(bson.M{"_id": dummyID}).Limit(1).Iter().All(&found)
	if err != nil {
		return nil, c.NewReadError("Failed to read item", err)
	}
	if len(found) == 0 {
		return nil, nil
	}

	return found[0], nil
}

func (c *DummyMongoDbDataAccess) CreateDummy(dummy *Dummy) (*Dummy, error) {
	if dummy.ID == "" {
		dummy.ID = c.CreateUUID()
	}

	err := c.collection.Insert(dummy)
	if err != nil {
		return nil, c.NewWriteError("Failed to insert item", err)
	}

	return dummy, nil
}

func (c *DummyMongoDbDataAccess) UpdateDummy(dummyID string, dummy *runtime.DynamicMap) (*Dummy, error) {
	if dummy == nil { return nil, nil }
    
    item, err := c.GetDummyById(dummyID)
    if err != nil { return nil, err }
    if item == nil { return item, nil }

    dummy.AssignTo(item)
    err = c.collection.UpdateId(dummyID, item)
    if err != nil {
        return nil, c.NewWriteError("Failed to update item", err)
    }

	return item, nil
}

func (c *DummyMongoDbDataAccess) DeleteDummy(dummyID string) error {
	_, err := c.collection.RemoveAll(bson.M{"_id": dummyID})
	if err != nil {
		return c.NewWriteError("Failed to remove item", err)
	}

	return nil
}
