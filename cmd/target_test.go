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

var _ = Describe("Commands", func() {
	var binLocation string
	var binFile string
	var StateFile string
	var target string
	var configFlag string
	var tempDir string
	var err error

	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s", runtime.GOOS)
		binFile = fmt.Sprintf("%s/workflow-cli", binLocation)

		tempDir, err = ioutil.TempDir("", "")
		Expect(err).To(BeNil())
		StateFile = fmt.Sprintf("%s/.cli", tempDir)
		configFlag = fmt.Sprintf("--config=%s", StateFile)

		port := 8080 + config.GinkgoConfig.ParallelNode

		if https {
			target = fmt.Sprintf("https://localhost:%d", port)
		} else {
			target = fmt.Sprintf("http://localhost:%d", port)
		}

	})

	AfterEach(func() {
		os.Remove(StateFile)
		os.RemoveAll(tempDir)
	})

	Describe("Test the commands", func() {

		Context("call target with valid input", func() {
			It("INTEGRATION should set info for the target", func() {
				// Set up command to test
				cmd := exec.Command(binFile, "target", target, configFlag)

				output, err := cmd.CombinedOutput()
				Expect(err).To(BeNil())

				// Verify state file is created
				_, err = ioutil.ReadFile(StateFile)
				Expect(err).ToNot(HaveOccurred())

				// Verify output
				expectedString := fmt.Sprintf("Target set to %s\n", target)

				Expect(string(output)).To(ContainSubstring(expectedString))
			})
			AfterEach(func() {
				os.Remove(StateFile)
			})
		})
		Context("Test HTTP/HTTPS mismatch error handling", func() {
			It("INTEGRATION should fail if client and server mismatch with https", func() {
				// Ensure client and server are not using same scheme
				if https {
					target = "http://localhost:8080"
				} else {
					target = "https://localhost:8080"
				}

				cmd := exec.Command(binFile, "target", target, configFlag)
				output, err := cmd.CombinedOutput()
				Expect(err).To(BeNil())

				// Verify state file is Not created
				_, err = ioutil.ReadFile(StateFile)
				Expect(err).To(HaveOccurred())

				// Verify output
				Expect(string(output)).To(ContainSubstring("Error"))
			})
			AfterEach(func() {
				os.Remove(StateFile)
			})
		})
		Context("Ater target has been set", func() {
			BeforeEach(func() {
				// Ensure file does not exist before test
				_, err := os.Stat(StateFile)
				Expect(os.IsNotExist(err))

				//Set the target
				cmd := exec.Command(binFile, "target", target, configFlag)
				_, err = cmd.CombinedOutput()

				Expect(err).To(BeNil())

				// Verify state file is created
				_, err = ioutil.ReadFile(StateFile)
				Expect(err).ToNot(HaveOccurred())

			})
			AfterEach(func() {
				_ = os.Remove(StateFile)
			})
			It("INTEGRATION should display endpoint after target", func() {
				cmd := exec.Command(binFile, "target", configFlag)
				output, err := cmd.CombinedOutput()
				Expect(err).To(BeNil())

				Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Current target is %s\n", target)))
			})
		})
		Context("no target has been set", func() {
			It("INTEGRATION should display no target set", func() {

				os.Remove(StateFile)

				cmd := exec.Command(binFile, "target", configFlag)
				output, err := cmd.CombinedOutput()
				Expect(err).To(BeNil())

				Expect(string(output)).To(ContainSubstring(fmt.Sprintf("No target set")))
			})
			AfterEach(func() {
				os.Remove(StateFile)
			})
		})
	})
})
