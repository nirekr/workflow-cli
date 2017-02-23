package main

import (
	"eos2git.cec.lab.emc.com/VCE-Symphony/workflow-cli/mock"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infof("Starting mock REST endpoint")
	mock.CreateMock()
}
