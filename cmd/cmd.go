package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/remote"
	"os"
	"strings"
)

// cmdCmd runs a single command against the development environment, and streams the results back.
// This opens an SSH connection, runs the command, and then closes the connection.
var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Run a command in the remove dev environment",
	Long:  `Run a single command (over SSH) within the remote dev environment.`,
	Run:   runCmdCommand,
}

func init() {
	// Allow users to pass any argument to `vessel cmd` without it
	// be interpreted as a flag for the `cmd` sub-command
	cmdCmd.Flags().SetInterspersed(false)
}

// runCmdCommand runs the command given within the development environment,
// streaming the output back to the client
func runCmdCommand(cmd *cobra.Command, args []string) {
	cfg, err := config.RetrieveProjectConfig(ConfigPath)

	if err != nil {
		logger.GetLogger().Error("command", "cmd", "error", err)
		PrintIfVerbose(Verbose, err, "could not read configuration")

		os.Exit(1)
	}

	connection := remote.NewConnection(&cfg.Remote)

	if err := connection.Cmd(strings.Join(args, " ")); err != nil {
		logger.GetLogger().Error("command", "cmd", "error", err)
		PrintIfVerbose(Verbose, err, "could not run given command")

		os.Exit(1)
	}
}
