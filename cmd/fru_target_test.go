package cmd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Commands", func() {
	var binLocation string
	var server *ghttp.Server
	var StateFile string

	BeforeEach(func() {
		binLocation = fmt.Sprintf("../bin/%s/workflow-cli", runtime.GOOS)
		server = ghttp.NewServer()
		dir, err := homedir.Dir()
		Expect(err).ToNot(HaveOccurred())
		StateFile = fmt.Sprintf("%s/.fru", dir)
	})

	AfterEach(func() {
		server.Close()
		os.Remove(StateFile)
	})

	Describe("Test the  commands", func() {
		Context("When calling 'target' with invalid url", func() {
			It("UNIT should check input with JSON Marshal and fail", func() {
				responseBytes := []byte("THIS SHOULD FAIL")
				expectedOutput := fmt.Sprintf("Invalid response: %s\n", responseBytes)

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/about"),
						ghttp.RespondWith(http.StatusOK, responseBytes),
					),
				)

				// Set up command to test
				cmd := exec.Command(binLocation, "fru", "target", server.URL())
				out, _ := cmd.StderrPipe()
				cmd.Start()

				// Capture Standard Output to verify
				buf := new(bytes.Buffer)
				buf.ReadFrom(out)
				s := buf.String()

				// Verify output
				Expect(s).To(ContainSubstring(expectedOutput))

			})
		})
		Context("call target with valid input", func() {
			It("INTEGRATION should set info for the target", func() {

				responseString := "up and running"
				expectedResponseData := []byte(responseString)
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/about"),
						ghttp.RespondWith(http.StatusOK, expectedResponseData),
					),
				)

				// Set up command to test
				cmd := exec.Command(binLocation, "fru", "target", server.URL())
				out, _ := cmd.StdoutPipe()
				cmd.Start()

				// Capture Standard Output to verify
				buf := new(bytes.Buffer)
				buf.ReadFrom(out)
				s := buf.String()
				fmt.Printf(s)

				// Verify state file is created
				_, err := ioutil.ReadFile(StateFile)

				// Verify output
				expectedString := fmt.Sprintf("Target set to %s\n", server.URL())

				Expect(s).To(Equal(expectedString))
				Expect(err).ToNot(HaveOccurred())
				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})
		})
		Context("Ater target has been set", func() {
			BeforeEach(func() {
				_, err := os.Stat(StateFile)
				Expect(os.IsNotExist(err))
				urlObj, err := url.Parse(server.URL())
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
				cmd := exec.Command(binLocation, "fru", "target")
				out, _ := cmd.StdoutPipe()
				cmd.Start()

				// Capture Standard Output to verify
				buf := new(bytes.Buffer)
				buf.ReadFrom(out)
				s := buf.String()

				Expect(s).To(Equal(fmt.Sprintf("Current target is %s\n", server.URL())))
			})
		})
		Context("no target has been set", func() {
			It("INTEGRATION should display no target set", func() {
				cmd := exec.Command(binLocation, "fru", "target")
				out, _ := cmd.StdoutPipe()
				cmd.Start()

				buf := new(bytes.Buffer)
				buf.ReadFrom(out)
				s := buf.String()

				Expect(s).To(Equal(fmt.Sprintf("No target set.\n")))
			})
		})
	})
})
