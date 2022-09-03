package cmd

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/mutagen"
	"os"
	"strings"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a remote dev session",
	Long:  `Connect to a remote server and start developing.`,
	Run:   runStartCommand,
}

func init() {
	startCmd.Flags().StringVarP(&ConfigPath, "config-file", "c", "vessel.yml", "Configuration file to read from")
}

// runStartCommand starts Mutagen sync and forwarding sessions based on
// configuration in the vessel.yml file. It allows you to start developing remotely!
func runStartCommand(cmd *cobra.Command, args []string) {
	cfg, err := config.RetrieveProjectConfig(ConfigPath)

	if err != nil {
		logger.GetLogger().Error("command", "start", "msg", "could not read configuration", "error", err)
		PrintIfVerbose(Verbose, err, "error starting syncing session")

		os.Exit(1)
	}

	// Get mutagen session name
	name := slug.Make("vessel-" + cfg.Name)

	// todo: We assume local path is "."
	_, err = mutagen.Sync(name, cfg.Remote.Alias, ".", cfg.Remote.RemotePath)

	if err != nil {
		logger.GetLogger().Error("command", "start", "msg", "error starting syncing session", "error", err)
		PrintIfVerbose(Verbose, err, "error starting syncing session")

		os.Exit(1)
	}

	// Forward multiple ports
	for k, p := range cfg.Forwarding {
		ports := strings.Split(p, ":")

		if len(ports) != 2 {
			logger.GetLogger().Error("command", "start", "msg", "invalid port forwarding configuration", "error", fmt.Errorf("port forwarding configuration must define both ports to forward separated by a colon, e.g. 8000:8000"))
			PrintIfVerbose(Verbose, err, "error starting syncing session")

			os.Exit(1)
		}

		_, err = mutagen.Forward(fmt.Sprintf("%s-%d", name, k), "tcp:127.0.0.1:"+ports[0], cfg.Remote.Alias, "tcp:127.0.0.1:"+ports[1])

		if err != nil {
			logger.GetLogger().Error("command", "start", "msg", "error starting forwarding session", "error", err)
			PrintIfVerbose(Verbose, err, "error starting forwarding session")

			os.Exit(1)
		}
	}

	fmt.Println("Started development session")
}
