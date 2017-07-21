//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/dellemc-symphony/workflow-cli/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// targetCmd represents the target command
var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "IP endpoint to target",
	Long: `IP endpoint to target.
This command will attempt a call to the IP specified, and if successful,
will store the IP for future use in the config file (default is ~/.cli)
Usage: workflow-cli fru target http://<ip address>:<port>
ex.: workflow-cli fru target http://192.168.1.1:80`,

	Run: func(c *cobra.Command, args []string) {

		// Argument required
		if len(args) > 1 {
			log.Warnf("Too Many Arguments.\n")
			return
		} else if len(args) < 1 {
			// Check and see if the endpoint has been set
			if _, err := os.Stat(targetFile); err == nil {
				fileContent, err := ioutil.ReadFile(targetFile)
				if err != nil {
					log.Fatal(err)
				}
				// Unmarshal data and print
				urlObject := url.URL{}
				err = json.Unmarshal(fileContent, &urlObject)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Current target is ")
				fmt.Printf("%s://%s\n", urlObject.Scheme, urlObject.Host)
			} else {
				log.Warnf("No target set: %s", err)
				return
			}
		} else {
			// Parse and validate argument
			targetURL, err := url.Parse(args[0])
			if err != nil {
				log.Warnf("Could not convert arg to IP Address: %s", err)
				return
			}

			if targetURL.Host == "" || targetURL.Scheme == "" {
				log.Warnf("Please enter a valid target url. ex: http://192.168.1.1:80\n")
				return
			}

			targetURL.Path = "about"
			_, err = utils.GetURL(*targetURL)
			if err != nil {
				log.Warnf(err.Error())
				return
			}
			targetURL.Path = ""

			fmt.Printf("Target set to %s://%s\n", targetURL.Scheme, targetURL.Host)

			// Store target URL of valid endpoint
			targetb, err := json.Marshal(targetURL)
			if err != nil {
				log.Warnf("Could not marshal IP Address to JSON: %s", err)
				return
			}

			err = ioutil.WriteFile(targetFile, targetb, 0666)
			if err != nil {
				log.Warnf("Error storing IP address to file: %s", err)
				return
			}
		}
		// Success
	},
}

func init() {
	RootCmd.AddCommand(targetCmd)

}
