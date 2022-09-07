package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"os"
)

var Version = "dev"

// Verbose sets the verbosity of error messaging for each command (show or hide specific error messages)
var Verbose bool

// ConfigPath is used by many commands to store the path to the vessel.yml file
var ConfigPath string

func PrintIfVerbose(verbose bool, err error, fallback string) {
	if verbose {
		fmt.Println(err)
	} else {
		fmt.Println(fallback)
	}
}

// rootCmd is the root/first command. All other commands are "under" this root command.
// The rootCmd is an alias for the "cmd" subcommand.
var rootCmd = &cobra.Command{
	Use:   "vessel",
	Short: "Vessel makes remote dev feel local",
	Long: `Remote development you'd swear was local.
  Find out more at https://vessel.app`,
	Run: runCmdCommand, // Alias for "vessel cmd ..."
}

// Execute registers all other commands. This is called by the main package.
func Execute() {
	commands := []*cobra.Command{
		authCmd,
		cmdCmd,
		initCmd,
		openCmd,
		sshCmd,
		startCmd,
		stopCmd,
	}

	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(commands...)
	if err := rootCmd.Execute(); err != nil {
		logger.GetLogger().Error("command", "root", "error", err)
		PrintIfVerbose(Verbose, err, "could not run given command")

		os.Exit(1)
	}
}
