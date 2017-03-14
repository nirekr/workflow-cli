//
// Copyright © 2017 Dell Inc. or its subsidiaries. All Rights Reserved.
// VCE Confidential/Proprietary Information
//

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var target string

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

	RootCmd.PersistentFlags().StringVar(&target, "target", "http://localhost:8080", "target gateway endpoint")
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "default", "config file (default is $HOME/.cli)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	//viper.SetConfigName(".cli")  // name of config file (without extension)
	//viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv() // read in environment variables that match

	if configFile == "default" {
		home := viper.Get("HOME")
		configFile = fmt.Sprintf("%s/.cli", home)

	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}