package cmd_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FruData", func() {
	var binLocation string
	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
	})

	Context("When command is called", func() {
		It("UNIT should print called message", func() {
			cmd := exec.Command(binLocation, "fru", "data")
			out, _ := cmd.StdoutPipe()
			cmd.Start()

			buf := new(bytes.Buffer)
			buf.ReadFrom(out)
			Expect(buf.String()).To(Equal("data called\n"))
		})
	})
})
