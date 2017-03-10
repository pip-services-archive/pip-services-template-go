package build

import (
    "github.com/pip-services/pip-services-runtime-go"    
    "github.com/pip-services/pip-services-runtime-go/build"    
    "github.com/pip-services/pip-services-template-go/db"    
    "github.com/pip-services/pip-services-template-go/logic"    
    "github.com/pip-services/pip-services-template-go/deps/version1"    
    "github.com/pip-services/pip-services-template-go/api/version1"    
)

func NewDummyBuilder() *build.Builder {
    types := runtime.NewMapAndSet(
        "db.file", func(c *runtime.DynamicMap) runtime.IComponent { return db.NewDummyFileDataAccess(c) },
        "db.mongodb", func(c *runtime.DynamicMap) runtime.IComponent { return db.NewDummyMongoDbDataAccess(c) },
        "deps.dummy.rest", func(c *runtime.DynamicMap) runtime.IComponent { return deps_v1.NewDummyRestClient(c) },
        "ctrl.default", func(c *runtime.DynamicMap) runtime.IComponent { return logic.NewDummyController(c) },
        "api.version1.rest", func(c *runtime.DynamicMap) runtime.IComponent { return api_v1.NewDummyRestService(c) },
        "api.default.rest", func(c *runtime.DynamicMap) runtime.IComponent { return api_v1.NewDummyRestService(c) } )
    return build.NewBuilder(types, "Dummy.Builder")
}