package main

import (
	"github.com/dellemc-symphony/workflow-cli/mock"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infof("Starting mock REST endpoint")
	mock.CreateMock()
}
