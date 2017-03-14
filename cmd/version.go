package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var binaryName, buildDate, commitHash, goVersion, releaseVersion string

// versioinCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints metadata version information about this CLI tool.",
	Long:  `workflow-cli version`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(binaryName)
		log.Println("  Release version: " + releaseVersion)
		log.Println("  Built On: " + buildDate)
		log.Println("  Commit Hash: " + commitHash)
		log.Println("  Go version: " + goVersion)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
