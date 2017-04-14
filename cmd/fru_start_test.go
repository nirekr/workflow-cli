package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/dellemc-symphony/workflow-cli/resources"
	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FruStart", func() {
	var binLocation string
	var endpointLocation string
	var StateFile string
	var target string

	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		endpointLocation = fmt.Sprintf("../bin/%s/endpoint.yaml", runtime.GOOS)

		dir, err := homedir.Dir()
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", dir)

		if https {
			target = "https://localhost:8080"
		} else {
			target = "http://localhost:8080"
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

	Context("When the endpoint file has all fields filled", func() {
		It("UNIT should run the 'fru start' without user input", func() {
			err := resources.WriteEndpointsFile("AllFields", endpointLocation)
			Expect(err).To(BeNil())

			cmd := exec.Command(binLocation, "fru", "start")

			output, err := cmd.CombinedOutput()
			Expect(err).To(BeNil())
			Expect(output).To(ContainSubstring("Workflow complete"))

		})
	})

	Context("When the endpoint file has endpoints but not credentials", func() {
		It("UNIT should run the 'fru start' and prompt for input", func() {
			err := resources.WriteEndpointsFile("MissingCredentials", endpointLocation)
			Expect(err).To(BeNil())

			cmd := exec.Command(binLocation, "fru", "start")

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

			io.WriteString(stdin, "a\n")
			time.Sleep(500 * time.Millisecond)

			io.WriteString(stdin, "b\n")
			time.Sleep(500 * time.Millisecond)

			io.WriteString(stdin, "1\n")
			time.Sleep(500 * time.Millisecond)

			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			buf := new(bytes.Buffer)
			buf.ReadFrom(stderr)
			Expect(buf.String()).To(ContainSubstring("Workflow complete"))

			cmd.Wait()

		})
	})

	Context("When the endpoint file is missing", func() {
		It("UNIT should prompt for endpoints and creds", func() {
			cmd := exec.Command(binLocation, "fru", "start")

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

			io.WriteString(stdin, "a\n")
			time.Sleep(500 * time.Millisecond)

			io.WriteString(stdin, "b\n")
			time.Sleep(500 * time.Millisecond)

			io.WriteString(stdin, "c\n")
			time.Sleep(500 * time.Millisecond)

			io.WriteString(stdin, "1\n")
			time.Sleep(500 * time.Millisecond)

			io.WriteString(stdin, "2\n")
			time.Sleep(500 * time.Millisecond)

			io.WriteString(stdin, "3\n")
			time.Sleep(500 * time.Millisecond)

			buf := new(bytes.Buffer)
			buf.ReadFrom(stderr)
			Expect(buf.String()).To(ContainSubstring("endpoint.yaml not found, will prompt user"))
			Expect(buf.String()).To(ContainSubstring("Workflow complete"))

		})
	})

})
