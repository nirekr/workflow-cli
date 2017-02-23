//
// Copyright © 2017 Dell Inc. or its subsidiaries. All Rights Reserved.
// VCE Confidential/Proprietary Information
//

package cmd

import (
	"eos2git.cec.lab.emc.com/VCE-Symphony/workflow-cli/frutaskrunner"
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

		r, err := frutaskrunner.InitiateWorkflow(target)
		if err != nil {
			log.Warnf("Error starting FRU task: %s", err)
		}

		err = frutaskrunner.RunTask(r)
		if err != nil {
			log.Warnf("Error running FRU task: %s", err)
		}
	},
}

func init() {
	fruCmd.AddCommand(startCmd)
}
