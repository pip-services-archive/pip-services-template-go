package deps_v1

import (
    "github.com/pip-services/pip-services-runtime-go"
    "github.com/pip-services/pip-services-template-go/db"
)

type IDummyClient interface {
    GetDummies(filter *runtime.FilterParams, paging *runtime.PagingParams) (*db.DummyDataPage, error)
    GetDummyById(dummyID string) (*db.Dummy, error)
    CreateDummy(dummy *db.Dummy) (*db.Dummy, error)
    UpdateDummy(dummyID string, dummy *runtime.DynamicMap) (*db.Dummy, error)
    DeleteDummy(dummyID string) error
}