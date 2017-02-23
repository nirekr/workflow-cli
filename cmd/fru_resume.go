//
// Copyright © 2017 Dell Inc. or its subsidiaries. All Rights Reserved.
// VCE Confidential/Proprietary Information
//

package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"eos2git.cec.lab.emc.com/VCE-Symphony/workflow-cli/frutaskrunner"
	"eos2git.cec.lab.emc.com/VCE-Symphony/workflow-cli/models"
	"eos2git.cec.lab.emc.com/VCE-Symphony/workflow-cli/utils"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
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

		// Do GET to /fru/api/workflow
		client := cleanhttp.DefaultClient()
		log.Printf("Server target is %s\n", target)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(target+"/fru/api/workflow"), nil)
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
		input := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Warnf("Error reading addr: %s", err)
		}

		selector, err := strconv.Atoi(input)
		if err != nil {
			log.Errorf("Error reading user input: %s", err)
		}
		uriSplit := strings.Split(workflows[selector-1].URI, "/")
		resumeID := uriSplit[len(uriSplit)-1]

		// Do GET to /fru/api/workflow/{task-id} to get response
		req1, err := http.NewRequest(http.MethodGet, fmt.Sprintf(target+"/fru/api/workflow/%s", resumeID), nil)
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

		//log.Infof("req1 is %+v\n", req1)
		//log.Infof("resp1 is %+v\n", resp1)
		//log.Infof("r is %+v\n", r)

		err = frutaskrunner.RunTask(r)
		if err != nil {
			log.Warnf("Error running FRU task: %s", err)
		}
	},
}

func init() {
	fruCmd.AddCommand(resumeCmd)
}
