package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

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
		startCmd,
		sshCmd,
		cmdCmd,
	}
	rootCmd.AddCommand(commands...)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
