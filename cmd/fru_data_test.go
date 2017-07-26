//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

var _ = Describe("FruData", func() {
	var binFile string
	var StateFile string
	var configFlag string
	var target string
	var tempDir string
	var err error

	BeforeEach(func() {
		binFile = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)

		tempDir, err = ioutil.TempDir("", "")
		Expect(err).To(BeNil())
		StateFile = fmt.Sprintf("%s/.cli", tempDir)
		configFlag = fmt.Sprintf("--config=%s", StateFile)

		nodeTestPort := 8080 + config.GinkgoConfig.ParallelNode
		if https {
			target = fmt.Sprintf("https://localhost:%d", nodeTestPort)
		} else {
			target = fmt.Sprintf("http://localhost:%d", nodeTestPort)
		}

		cmd := exec.Command(binFile, "target", target, configFlag)
		_, err = cmd.CombinedOutput()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		os.Remove(StateFile)
		os.RemoveAll(tempDir)
	})

	Context("When command is called", func() {
		It("UNIT should print called message", func() {
			cmd := exec.Command(binFile, "fru", "data", "47daaa4d-8c4f-40cd-84db-901963d1fc0c", configFlag)
			output, err := cmd.CombinedOutput()
			Expect(err).To(BeNil())

			Expect(string(output)).To(ContainSubstring("Fru Data"))

		})
	})
})
