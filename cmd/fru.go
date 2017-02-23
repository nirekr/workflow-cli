//
// Copyright © 2017 Dell Inc. or its subsidiaries. All Rights Reserved.
// VCE Confidential/Proprietary Information
//

package cmd

import "github.com/spf13/cobra"

// fruCmd represents the fru command
var fruCmd = &cobra.Command{
	Use:   "fru",
	Short: "fru commands control field replacement and debugging of FRUs",
	Long: `Use the commands below to execute FRU replacement workflows:

Ensure you have the Customer vCenter IP/Hostname and Credentials along with
the Customer ScaleIO Gateway IP and credentials.`,
}

func init() {
	RootCmd.AddCommand(fruCmd)
	// forces a password auth on every command
	// need to cache this somehow
	// auth.UserAuth() // TEMPORARILY DISABLED UNTIL SECURITY CONVERSATIONS
}
