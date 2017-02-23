package cmd_test

import (
	"eos2git.cec.lab.emc.com/VCE-Symphony/workflow-cli/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/gin-gonic/gin.v1"

	"testing"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

var router *gin.Engine

var _ = BeforeSuite(func() {
	go mock.CreateMock()
})

var _ = AfterSuite(func() {
	mock.StopMock()
})
