//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/dellemc-symphony/workflow-cli/frutaskrunner"
	"github.com/dellemc-symphony/workflow-cli/models"
	"github.com/dellemc-symphony/workflow-cli/transport"
	"github.com/dellemc-symphony/workflow-cli/utils"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "This command will resume a failed or paused VxRack FRU replacement operation.",
	Long: `This will provide a list of running tasks or automatically choose the currently
running task if there is only 1.

Use this command to resume a failed or stopped workflow operation.`,
	Run: func(cmd *cobra.Command, args []string) {

		fileContent, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}
		// Unmarshal data and print
		urlObject := url.URL{}
		err = json.Unmarshal(fileContent, &urlObject)
		if err != nil {
			log.Fatal(err)
		}

		// Do GET to /fru/api/workflow
		client, err := transport.NewClient(urlObject.String())
		if err != nil {
			log.Fatalf(err.Error())
		}

		log.Printf("Server target is %s\n", urlObject.String())
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlObject.String()+"/fru/api/workflow"), nil)
		if err != nil {
			log.Warnf("%s", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if resp.StatusCode < 200 || 300 < resp.StatusCode {
			log.Warnf("Non-Success status(%d): %s", resp.StatusCode, resp.Status)
		}

		workflows := models.Workflows{}
		if err = utils.DecodeBody(resp, &workflows); err != nil {
			log.Warnf("Decoding Response: %s", err)
		}

		// Print result
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"WorkflowID", "Select"})

		for index, workflow := range workflows {
			uriSplit := strings.Split(workflow.URI, "/")
			workflowID := uriSplit[len(uriSplit)-1]
			table.Append([]string{workflowID, fmt.Sprintf("%d", index+1)})
		}
		table.Render()

		// Ask which task to resume, by task-id
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Please enter resume task: ")

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Warnf("Error reading addr: %s", err)
		}

		input := scanner.Text()

		selector, err := strconv.Atoi(input)
		if err != nil {
			log.Errorf("Error reading user input: %s", err)
		}
		uriSplit := strings.Split(workflows[selector-1].URI, "/")
		resumeID := uriSplit[len(uriSplit)-1]

		// Do GET to /fru/api/workflow/{task-id} to get response
		req1, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlObject.String()+"/fru/api/workflow/%s", resumeID), nil)
		if err != nil {
			log.Warnf("%s", err)
		}

		resp1, err := client.Do(req1)
		if err != nil {
			log.Fatalf("%s", err)
		}

		if resp1.StatusCode < 200 || 300 < resp1.StatusCode {
			log.Warnf("Tried (%s) %s", req1.Method, req1.URL.String())
			log.Warnf("Non-Success status(%d): %s", resp1.StatusCode, resp1.Status)
		}

		// pass response to response runner
		r := models.Response{}
		if err = utils.DecodeBody(resp1, &r); err != nil {
			log.Warnf("Decoding Response: %s", err)
		}

		err = frutaskrunner.RunTask(r, urlObject.String())
		if err != nil {
			log.Warnf("Error running FRU task: %s", err)
		}
	},
}

func init() {
	fruCmd.AddCommand(resumeCmd)
}
