package cmd

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/mutagen"
	"os"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a remote dev session",
	Long:  `Disconnect from an active remote server.`,
	Run:   runStopCommand,
}

func init() {
	stopCmd.Flags().StringVarP(&ConfigPath, "config-file", "c", "vessel.yml", "Configuration file to read from")
}

// runStartCommand starts Mutagen sync and forwarding sessions based on
// configuration in the vessel.yml file. It allows you to start developing remotely!
func runStopCommand(cmd *cobra.Command, args []string) {
	cfg, err := config.RetrieveProjectConfig(ConfigPath)

	if err != nil {
		logger.GetLogger().Error("command", "start", "msg", "error stopping development session", "error", err)
		PrintIfVerbose(Verbose, err, "error stopping development session")

		os.Exit(1)
	}

	// Generate mutagen session name
	// Get mutagen session name
	name := slug.Make("vessel-" + cfg.Name)

	err = mutagen.StopSession(name)

	if err != nil {
		logger.GetLogger().Error("command", "start", "msg", "error stopping development session", "error", err)
		PrintIfVerbose(Verbose, err, "error stopping development session")

		os.Exit(1)
	}

	fmt.Println("Stopped development session")
}
