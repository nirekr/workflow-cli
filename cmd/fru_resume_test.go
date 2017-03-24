package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/dellemc-symphony/workflow-cli/frutaskrunner"
	homedir "github.com/mitchellh/go-homedir"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FruResume", func() {

	var binLocation string
	var target string
	var StateFile string
	BeforeEach(func() {
		dir, err := homedir.Dir()
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", dir)

		binLocation = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		if https {
			target = "https://localhost:8080"
		} else {
			target = "http://localhost:8080"
		}

		cmd := exec.Command(binLocation, "target", target)
		err = cmd.Run()
		Expect(err).To(BeNil())

	})
	AfterEach(func() {
		os.Remove(StateFile)
	})

	Context("When command is called", func() {
		It("UNIT should print called message", func() {
			_, err := frutaskrunner.InitiateWorkflow(target)
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
			io.WriteString(stdin, "1\n")
			time.Sleep(500 * time.Millisecond)

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

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stderr)
			Expect(errBuf.String()).To(ContainSubstring("Workflow complete"))

			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdout)
			tableString := `+----------------------+--------+
|      WORKFLOWID      | SELECT |
+----------------------+--------+
| 123abc-456def-789ghi |      1 |
+----------------------+--------+
`
			Expect(outBuf.String()).To(ContainSubstring(tableString))

			cmd.Wait()

		})
	})
})
