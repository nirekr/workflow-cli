//
// Copyright © 2017 Dell Inc. or its subsidiaries. All Rights Reserved.
// VCE Confidential/Proprietary Information
//

package cmd

import (
	"fmt"

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
		// TODO: add the API call to data collection
		fmt.Println("data called")
	},
}

func init() {
	fruCmd.AddCommand(dataCmd)
}
