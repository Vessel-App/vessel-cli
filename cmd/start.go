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
	"strings"
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

		// stop syncing
		errSync := mutagen.StopSync(name)

		// stop forwarding
		errForward := mutagen.StopForward(name)

		if errSync != nil || errForward != nil {
			fmt.Println("Error disconnecting from development server")
			os.Exit(1)
		}

		done <- true
	}()

	// Pause until user decides we are done
	<-done
}
