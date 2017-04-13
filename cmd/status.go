//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/dellemc-symphony/workflow-cli/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Provides status call to the system",
	Long: `Only to be used to determine that the system is
up and running. Does not provide information about VxRack system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		urlObject.Path = "about"
		resp, err := utils.GetURL(urlObject)
		if err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Printf("Status: \n%s\n", resp)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
