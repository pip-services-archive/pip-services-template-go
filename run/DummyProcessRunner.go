package run

import (
    "github.com/pip-services/pip-services-runtime-go/run"
    "github.com/pip-services/pip-services-template-go/build"        
)

type DummyProcessRunner struct {
    run.ProcessRunner    
}

func NewDummyProcessRunner() *DummyProcessRunner {
    builder := build.NewDummyBuilder()
    return &DummyProcessRunner { ProcessRunner: *run.NewProcessRunner(builder) }
}
