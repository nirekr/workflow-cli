//
// Copyright © 2017 Dell Inc. or its subsidiaries. All Rights Reserved.
// VCE Confidential/Proprietary Information
//

package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Provides status call to the system",
	Long: `Only to be used to determine that the system is
up and running. Does not provide information about VxRack system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// this is a mock response call
		//var addrs string
		var ip net.IP
		interfaces, err := net.Interfaces()
		if err != nil {
			return err
		}
		for _, val := range interfaces {
			addr, err := val.Addrs()
			if err != nil {
				return err
			}
			for _, a := range addr {
				switch v := a.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
			}
		}
		fmt.Printf("Gateway is UP on %s\n", ip)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
