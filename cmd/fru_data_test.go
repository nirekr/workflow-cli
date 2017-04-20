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
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FruData", func() {
	var binLocation string
	var StateFile string
	var target string
	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		dir, err := homedir.Dir()
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", dir)

		if https {
			target = "https://localhost:8080"
		} else {
			target = "http://localhost:8080"
		}

		cmd := exec.Command(binLocation, "target", target)
		_, err = cmd.CombinedOutput()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		os.Remove(StateFile)
	})

	Context("When command is called", func() {
		It("UNIT should print called message", func() {
			cmd := exec.Command(binLocation, "fru", "data", "47daaa4d-8c4f-40cd-84db-901963d1fc0c")
			output, err := cmd.CombinedOutput()
			Expect(err).To(BeNil())

			Expect(string(output)).To(ContainSubstring("Fru Data"))

		})
	})
})
