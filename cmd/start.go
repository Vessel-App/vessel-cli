package cmd

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/mutagen"
	"os"
	"os/signal"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a remote dev session",
	Long:  `Connect to a remote server and start developing.`,
	Run:   runStartCommand,
}

var detach = false

func init() {
	startCmd.Flags().StringVarP(&ConfigPath, "config-file", "c", "vessel.yml", "Configuration file to read from")
	startCmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run in the background. Run `vessel stop` to stop the development session.")
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

	// Attempt to stop any currently running session to prevent duplicates
	// Note that we ignore errors
	mutagen.StopSession(name)

	// todo: We assume local path is "."
	err = mutagen.StartSession(name, ".", cfg)

	if err != nil {
		logger.GetLogger().Error("command", "start", "msg", "error starting syncing session", "error", err)
		PrintIfVerbose(Verbose, err, "error starting development session, attempting to cleanup session")

		// Cleanup any syncing or forwarding started
		// Note that we ignore errors
		mutagen.StopSession(name)

		os.Exit(1)
	}

	fmt.Println("Started development session")

	// If -d, --detach is used, we end here
	if detach {
		return
	}

	fmt.Println("Use crtl+c to stop the session")

	// Else we treat the command as long-running. We listen of os.Interrupt or os.Kill signals
	// (which work on Windows/Linux as per https://stackoverflow.com/a/35683558/1412984) and clean up
	// if those signals are received
	sigc := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	go func() {
		_ = <-sigc

		fmt.Println("\nStopping development session")

		err = mutagen.StopSession(name)

		if err != nil {
			fmt.Println("error stopping development session")
			os.Exit(1)
		}

		done <- true
	}()

	// Pause until user decides we are done
	<-done
}
