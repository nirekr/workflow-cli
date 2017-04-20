//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/url"

	"github.com/dellemc-symphony/workflow-cli/frutaskrunner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Executes the FRU replacement workflow with the Symphony FRU PAQX",
	Long: `This command will execute the VxRack FRU replacement operation.

The workflow will walk you through the process and allow you to start/stop at each step
as needed. Using the 'resume' command will allow a failed run to be restarted where it
left off`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Initiating workflow: quanta-replacement-d51b-esxi")

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

		r, err := frutaskrunner.InitiateWorkflow(urlObject.String())
		if err != nil {
			log.Warnf("Error starting FRU task: %s", err)
		}

		err = frutaskrunner.RunTask(r, urlObject.String())
		if err != nil {
			log.Warnf("Error running FRU task: %s", err)
		}
	},
}

func init() {
	fruCmd.AddCommand(startCmd)
}
