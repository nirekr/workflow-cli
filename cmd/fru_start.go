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
	"net/url"
	"os"
	"strings"

	"github.com/dellemc-symphony/workflow-cli/frutaskrunner"
	"github.com/dellemc-symphony/workflow-cli/resources"
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
		//Ask the user if they want the Endpoint file
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Do you want to use the 'endpoint.yaml'? [Y/N]: ")
		scanner.Scan()
		fmt.Printf("\n")
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading user input: %s", err)
		}

		input := scanner.Text()

		if strings.ToLower(input) == "y" {
			resources.UseEndpointFile = true
		} else {
			resources.UseEndpointFile = false
		}

		if resources.UseEndpointFile {
			// If they want the Endpoint file, validate it
			endpoints := resources.ParseEndpointsFile(resources.UseEndpointFile)
			for endpointName, endpoint := range endpoints {

				if endpoint.EndpointURL == "" {
					log.Fatalf("Endpoint file has invalid Endpoint entry for %s", endpointName)
				}
				if endpoint.Username == "" {
					log.Fatalf("Endpoint file has invalid Username entry for %s", endpointName)
				}
				if endpoint.Password == "" {
					log.Fatalf("Endpoint file has invalid Password entry for %s", endpointName)
				}

			}
		}

		log.Println("Initiating workflow: quanta-replacement-d51b-esxi")

		fileContent, err := ioutil.ReadFile(targetFile)
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
