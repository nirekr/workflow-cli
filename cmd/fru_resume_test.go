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

	"github.com/dellemc-symphony/workflow-cli/frutaskrunner"
	"github.com/dellemc-symphony/workflow-cli/resources"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/ginkgo/config"
)

var _ = Describe("FruResume", func() {

	var binLocation string
	var endpointLocation string
	var target string
	var StateFile string

	BeforeEach(func() {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", dir)

		binLocation = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		endpointLocation = fmt.Sprintf("../bin/%s/endpoint.yaml", runtime.GOOS)

                nodeTestPort := 8080 + config.GinkgoConfig.ParallelNode
                if https {
                        target = fmt.Sprintf("https://localhost:%d", nodeTestPort)
                } else {
                        target = fmt.Sprintf("http://localhost:%d", nodeTestPort)
                }

		cmd := exec.Command(binLocation, "target", target)
		err = cmd.Run()
		Expect(err).To(BeNil())

		// Remove any endpoint file previously in bin
		os.Remove(endpointLocation)
	})
	AfterEach(func() {
		os.Remove(StateFile)
	})

	Context("When command is called", func() {
		It("UNIT should print called message", func() {
			tableString := `+----------------------+--------+
|      WORKFLOWID      | SELECT |
+----------------------+--------+
| 123abc-456def-789ghi |      1 |
+----------------------+--------+
`

			err := resources.WriteEndpointsFile("AllFields", endpointLocation)
			Expect(err).To(BeNil())

			_, err = frutaskrunner.InitiateWorkflow(target)
			Expect(err).To(BeNil())

			cmd := exec.Command(binLocation, "fru", "resume")

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

			// Select workflow 1 to resume
			io.WriteString(stdin, "1\n")
			time.Sleep(500 * time.Millisecond)

			// Select node 2 to remove
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			// Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(35000 * time.Millisecond)

			//CONTINUE to allow node addition
			io.WriteString(stdin, "CONTINUE\n")
			time.Sleep(500 * time.Millisecond)

			// Select node 2 to add
			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			// Confirm selection
			io.WriteString(stdin, "Y\n")
			time.Sleep(500 * time.Millisecond)

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)

			Expect(outBuf.String()).To(ContainSubstring(tableString))
		})
	})
})
