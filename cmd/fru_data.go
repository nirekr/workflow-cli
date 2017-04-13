//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/dellemc-symphony/workflow-cli/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "This command will show data collected by the VxRack FRU replacement operation.",
	Long: `During the FRU replacement workflow, it will collect VxRack data relevant for
debugging and system visibility.

This command gives that data in a tabular format which can be used by the user to debug or validate
system configurations. It currently does not allow changes to the collected data sets.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Convert argument to REST call
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

		if len(args) != 0 {
			urlObject.Path = fmt.Sprintf("data/%s", args[0])
		} else {
			log.Fatalf("Error: Please provide a taskID!")
		}

		resp, err := utils.GetURL(urlObject)
		if err != nil {
			log.Fatalf(err.Error())
		}

		var data bytes.Buffer
		err = json.Indent(&data, resp.([]byte), "", "    ")
		if err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Printf("Fru Data: \n%s\n", data.Bytes())
	},
}

func init() {
	fruCmd.AddCommand(dataCmd)
}
