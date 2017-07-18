//
// Copyright (c) 2017 Dell Inc. or its subsidiaries.  All Rights Reserved.
// Dell EMC Confidential/Proprietary Information
//
//

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

//var target string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "workflow-cli",
	Short: "This is the Workflow CLI",
	Long: `This CLI is used to interact with the different PAQX
services.

Use this CLI to debug and run various commands against the system.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "default", "Location of the configuration file")
}

func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	if configFile == "default" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatalf("Could not read application directory: %s", err.Error())
		}
		configFile = fmt.Sprintf("%s/.cli", dir)

	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
