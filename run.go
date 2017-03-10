package main

import (
    "github.com/pip-services/pip-services-template-go/run"        
)

func main() {
    runner := run.NewDummyProcessRunner()
    runner.RunWithDefaultConfigFile("config.json")
}