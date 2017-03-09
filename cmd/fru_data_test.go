package cmd_test

import (
	"bytes"
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
	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		dir, err := homedir.Dir()
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.cli", dir)
	})

	AfterEach(func() {
		os.Remove(StateFile)
	})

	Context("When command is called", func() {
		It("UNIT should print called message", func() {
			cmd := exec.Command(binLocation, "target", "http://localhost:8080")
			cmd.Run()

			cmd = exec.Command(binLocation, "fru", "data", "47daaa4d-8c4f-40cd-84db-901963d1fc0c")
			out, _ := cmd.StdoutPipe()
			cmd.Start()

			buf := new(bytes.Buffer)
			buf.ReadFrom(out)
			Expect(buf.String()).To(ContainSubstring("YooHoo!"))

			cmd.Wait()
		})
	})
})
