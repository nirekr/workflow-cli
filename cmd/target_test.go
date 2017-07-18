//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commands", func() {
	var binLocation string
	var binFile string
	var StateFile string
	var target string

	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s", runtime.GOOS)
		binFile = fmt.Sprintf("%s/workflow-cli", binLocation)

		// configDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		// Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", binLocation)

		if https {
			target = "https://localhost:8080"
		} else {
			target = "http://localhost:8080"
		}

	})

	AfterEach(func() {
		os.Remove(StateFile)
	})

	Describe("Test the commands", func() {

		Context("call target with valid input", func() {
			It("INTEGRATION should set info for the target", func() {
				// Set up command to test
				cmd := exec.Command(binFile, "target", target)

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

				cmd := exec.Command(binFile, "target", target)
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
				_, err := os.Stat(StateFile)
				Expect(os.IsNotExist(err))
				urlObj, err := url.Parse(target)
				Expect(err).ToNot(HaveOccurred())
				urlBytes, err := json.Marshal(urlObj)
				Expect(err).ToNot(HaveOccurred())
				err = ioutil.WriteFile(StateFile, urlBytes, 0666)
				Expect(err).ToNot(HaveOccurred())
			})
			AfterEach(func() {
				_ = os.Remove(StateFile)
			})
			It("INTEGRATION should display endpoint after target", func() {
				cmd := exec.Command(binFile, "target")
				output, err := cmd.CombinedOutput()
				Expect(err).To(BeNil())

				Expect(string(output)).To(ContainSubstring(fmt.Sprintf("Current target is %s\n", target)))
			})
		})
		Context("no target has been set", func() {
			It("INTEGRATION should display no target set", func() {

				_, err := os.Stat(StateFile)
				if err == nil {
					// err == nil means the file exists, so we remove it
					os.Remove(StateFile)
				}

				cmd := exec.Command(binFile, "target")
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
