//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/ginkgo/config"
)

var _ = Describe("FruData", func() {
	var binFile string
	var StateFile string
	var target string
	BeforeEach(func() {
		binFile = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", dir)

                nodeTestPort := 8080 + config.GinkgoConfig.ParallelNode
                if https {
                        target = fmt.Sprintf("https://localhost:%d", nodeTestPort)
                } else {
                        target = fmt.Sprintf("http://localhost:%d", nodeTestPort)
                }

		cmd := exec.Command(binFile, "target", target)
		_, err = cmd.CombinedOutput()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		os.Remove(StateFile)
	})

	Context("When command is called", func() {
		It("UNIT should print called message", func() {
			cmd := exec.Command(binFile, "fru", "data", "47daaa4d-8c4f-40cd-84db-901963d1fc0c")
			output, err := cmd.CombinedOutput()
			Expect(err).To(BeNil())

			Expect(string(output)).To(ContainSubstring("Fru Data"))

		})
	})
})
