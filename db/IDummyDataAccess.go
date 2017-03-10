package db

import (
    "github.com/pip-services/pip-services-runtime-go"
)

type IDummyDataAccess interface {
    GetDummies(filter *runtime.FilterParams, paging *runtime.PagingParams) (*DummyDataPage, error)
    GetDummyById(dummyID string) (*Dummy, error)
    CreateDummy(dummy *Dummy) (*Dummy, error)
    UpdateDummy(dummyID string, dummy *runtime.DynamicMap) (*Dummy, error)
    DeleteDummy(dummyID string) error
}