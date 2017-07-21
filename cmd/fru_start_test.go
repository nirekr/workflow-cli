//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dellemc-symphony/workflow-cli/resources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FruStart", func() {
	var binFile string
	var endpointLocation string
	var StateFile string
	var target string
	var nodeList string
	var nodeSelection string
	var longDelay time.Duration

	BeforeEach(func() {

		longDelay = 35000

		binFile = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		endpointLocation = fmt.Sprintf("../bin/%s/endpoint.yaml", runtime.GOOS)

		nodeList = `+--------+----------+-------------+------------+--------------+------------------+
| SELECT | HOSTNAME | SERVICE TAG |  MGMT IP   | POWER STATUS | CONNECTION STATE |
+--------+----------+-------------+------------+--------------+------------------+
|      1 | node01   |     1234567 | 10.10.10.1 | poweredOn    | Connected        |
|      2 | node02   |    98765432 | 10.10.10.2 | poweredOn    | Connected        |
|      3 | node03   |    91827465 | 10.10.10.3 | poweredOn    | Connected        |
+--------+----------+-------------+------------+--------------+------------------+
`
		nodeSelection = `+----------+-------------+------------+--------------+------------------+
| HOSTNAME | SERVICE TAG |  MGMT IP   | POWER STATUS | CONNECTION STATE |
+----------+-------------+------------+--------------+------------------+
| node02   |    98765432 | 10.10.10.2 | poweredOn    | Connected        |
+----------+-------------+------------+--------------+------------------+
`

		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", dir)

		if https {
			target = "https://localhost:8080"
		} else {
			target = "http://localhost:8080"
		}

		cmd := exec.Command(binFile, "target", target)
		err = cmd.Run()
		Expect(err).To(BeNil())

		// Remove any endpoint file previously in bin
		os.Remove(endpointLocation)
	})
	AfterEach(func() {
		os.Remove(StateFile)
	})

	Context("When the endpoint file has all fields filled", func() {
		It("UNIT should run the 'fru start' without user input", func() {
			err := resources.WriteEndpointsFile("AllFields", endpointLocation)
			Expect(err).To(BeNil())

			startTime := time.Now()

			cmd := exec.Command(binFile, "fru", "start")

			stdin, err := cmd.StdinPipe()
			Expect(err).To(BeNil())
			defer stdin.Close()

			stdout, err := cmd.StdoutPipe()
			Expect(err).To(BeNil())
			defer stdout.Close()

			stderr, err := cmd.StderrPipe()
			Expect(err).To(BeNil())
			defer stderr.Close()

			cmd.Start()

			// Select y to use Endpoint.yaml
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			// Select node 2 for removal
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(longDelay * time.Millisecond)

			//CONTINUE to allow node addition
			io.WriteString(stdin, "CONTINUE\n")
			time.Sleep(500 * time.Millisecond)

			// Select node 2 for add
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(err).To(BeNil())
			Expect(outBuf.String()).To(ContainSubstring("It is now safe to remove the failed node from the rack."))
			Expect(outBuf.String()).To(ContainSubstring("When the new node has been racked and cabled, please power it on and type CONTINUE ..."))
			Expect(outBuf.String()).To(ContainSubstring(nodeList))
			Expect(outBuf.String()).To(ContainSubstring(nodeSelection))
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

			elapsedTime := time.Since(startTime)
			Expect(elapsedTime.Seconds()).To(BeNumerically(">", 5))
		})
	})

	Context("When the endpoint file has endpoints but not credentials", func() {
		It("UNIT should fail to parse the file", func() {
			err := resources.WriteEndpointsFile("MissingCredentials", endpointLocation)
			Expect(err).To(BeNil())

			cmd := exec.Command(binFile, "fru", "start")

			stdin, err := cmd.StdinPipe()
			Expect(err).To(BeNil())
			defer stdin.Close()

			stdout, err := cmd.StdoutPipe()
			Expect(err).To(BeNil())
			defer stdout.Close()

			stderr, err := cmd.StderrPipe()
			Expect(err).To(BeNil())
			defer stderr.Close()

			cmd.Start()

			// Select y to use Endpoint.yaml
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(errBuf.String()).To(ContainSubstring("Endpoint file has invalid Username entry"))

			cmd.Wait()

		})
	})

	Context("When the user enters mismatched passwords", func() {
		It("UNIT prompt for passwords again.", func() {
			os.Remove(StateFile)

			cmd := exec.Command(binFile, "fru", "start")

			stdin, err := cmd.StdinPipe()
			Expect(err).To(BeNil())
			defer stdin.Close()

			stdout, err := cmd.StdoutPipe()
			Expect(err).To(BeNil())
			defer stdout.Close()

			stderr, err := cmd.StderrPipe()
			Expect(err).To(BeNil())
			defer stderr.Close()

			cmd.Start()

			// Select y to use Endpoint.yaml
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for RackHD
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for RackHD
			io.WriteString(stdin, "RackHDUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for RackHD
			io.WriteString(stdin, "RackHDPassword_different\n")
			time.Sleep(500 * time.Millisecond)

			//Retry Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Retry Confirm Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for HostBMC
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for HostBMC
			io.WriteString(stdin, "HostBMCusername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for HostBMC
			io.WriteString(stdin, "HostBMCpassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for HostBMC
			io.WriteString(stdin, "HostBMCpassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for vCenter
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for vCenter
			io.WriteString(stdin, "vCenterUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for vCenter
			io.WriteString(stdin, "vCenterPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for vCenter
			io.WriteString(stdin, "vCenterPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for ScaleIOGateway
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for removal
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			// Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(longDelay * time.Millisecond)

			//CONTINUE to allow node addition
			io.WriteString(stdin, "CONTINUE\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for addition
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(outBuf.String()).To(ContainSubstring(nodeList))
			Expect(outBuf.String()).To(ContainSubstring(nodeSelection))
			Expect(outBuf.String()).To(ContainSubstring("Confirm rackhd Password"))
			Expect(errBuf.String()).To(ContainSubstring("Passwords for rackhd don't match"))
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

			cmd.Wait()

		})
	})

	Context("When the user enters an invalid endpoint", func() {
		It("UNIT prompt for endpoint again.", func() {
			os.Remove(StateFile)

			cmd := exec.Command(binFile, "fru", "start")

			stdin, err := cmd.StdinPipe()
			Expect(err).To(BeNil())
			defer stdin.Close()

			stdout, err := cmd.StdoutPipe()
			Expect(err).To(BeNil())
			defer stdout.Close()

			stderr, err := cmd.StderrPipe()
			Expect(err).To(BeNil())
			defer stderr.Close()

			cmd.Start()

			// Select y to use Endpoint.yaml
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for RackHD
			io.WriteString(stdin, "  http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for RackHD
			io.WriteString(stdin, "invalid-endpoint.com\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for RackHD
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for RackHD
			io.WriteString(stdin, "RackHDUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for RackHD
			io.WriteString(stdin, "RackHDPassword_different\n")
			time.Sleep(500 * time.Millisecond)

			//Retry Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Retry Confirm Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for HostBMC
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for HostBMC
			io.WriteString(stdin, "HostBMCusername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for HostBMC
			io.WriteString(stdin, "HostBMCpassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for HostBMC
			io.WriteString(stdin, "HostBMCpassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for vCenter
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for vCenter
			io.WriteString(stdin, "vCenterUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for vCenter
			io.WriteString(stdin, "vCenterPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for vCenter
			io.WriteString(stdin, "vCenterPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for ScaleIOGateway
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for removal
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			// Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(longDelay * time.Millisecond)

			//CONTINUE to allow node addition
			io.WriteString(stdin, "CONTINUE\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for addition
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(outBuf.String()).To(ContainSubstring(nodeList))
			Expect(outBuf.String()).To(ContainSubstring(nodeSelection))
			Expect(errBuf.String()).To(ContainSubstring("Endpoint must begin with http:// or https://"))
			Expect(errBuf.String()).To(ContainSubstring("Invalid URL"))
			Expect(outBuf.String()).To(ContainSubstring("Confirm rackhd Password"))
			Expect(errBuf.String()).To(ContainSubstring("Passwords for rackhd don't match"))
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

			cmd.Wait()

		})
	})

	Context("When the endpoint file is missing", func() {
		It("UNIT should prompt for endpoints and creds", func() {
			cmd := exec.Command(binFile, "fru", "start")

			stdin, err := cmd.StdinPipe()
			Expect(err).To(BeNil())
			defer stdin.Close()

			stdout, err := cmd.StdoutPipe()
			Expect(err).To(BeNil())
			defer stdout.Close()

			stderr, err := cmd.StderrPipe()
			Expect(err).To(BeNil())
			defer stderr.Close()

			cmd.Start()

			// Select N to not use Endpoint.yaml
			io.WriteString(stdin, "N\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for RackHD
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for RackHD
			io.WriteString(stdin, "RackHDUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for RackHD
			io.WriteString(stdin, "RackHDPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for HostBMC
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for HostBMC
			io.WriteString(stdin, "HostBMCusername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for HostBMC
			io.WriteString(stdin, "HostBMCpassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for HostBMC
			io.WriteString(stdin, "HostBMCpassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for vCenter
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for vCenter
			io.WriteString(stdin, "vCenterUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for vCenter
			io.WriteString(stdin, "vCenterPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for vCenter
			io.WriteString(stdin, "vCenterPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Endpoint for ScaleIOGateway
			io.WriteString(stdin, "http://10.10.10.10:8080\n")
			time.Sleep(500 * time.Millisecond)

			//Username for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayUsername\n")
			time.Sleep(500 * time.Millisecond)

			//Password for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm Password for ScaleIOGateway
			io.WriteString(stdin, "ScaleIOGatewayPassword\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for removal
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			// Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(longDelay * time.Millisecond)

			//CONTINUE to allow node addition
			io.WriteString(stdin, "CONTINUE\n")
			time.Sleep(500 * time.Millisecond)

			//Select node 2 for addition
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			//Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(outBuf.String()).To(ContainSubstring(nodeList))
			Expect(outBuf.String()).To(ContainSubstring(nodeSelection))
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

		})
	})

})
