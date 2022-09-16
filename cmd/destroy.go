package cmd

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/fly"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/mutagen"
	"github.com/vessel-app/vessel-cli/internal/util"
	"os"
	"time"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a dev environment",
	Long:  `Delete the development environment and related local Vessel files`,
	Run:   runDestroyCommand,
}

var localFiles bool
var shutUp bool

func init() {
	destroyCmd.Flags().BoolVarP(&localFiles, "files-only", "f", false, "Only delete local files, not the virtual machine")
	destroyCmd.Flags().BoolVarP(&shutUp, "quit", "q", false, "Delete without prompting for approval")
	destroyCmd.Flags().StringVarP(&ConfigPath, "config-file", "c", "vessel.yml", "Configuration file to read from")
}

func runDestroyCommand(cmd *cobra.Command, args []string) {
	cfg, err := config.RetrieveProjectConfig(ConfigPath)

	if err != nil {
		logger.GetLogger().Error("command", "destroy", "msg", "could not read configuration", "error", err)
		PrintIfVerbose(Verbose, err, "error reading project configuration file")

		os.Exit(1)
	}

	auth, err := config.RetrieveVesselConfig()

	if err != nil {
		logger.GetLogger().Error("command", "destroy", "msg", "could not get Fly API token from vessel config", "error", err)
		PrintIfVerbose(Verbose, err, "error retrieving Fly API token")

		os.Exit(1)
	}

	// Get mutagen session name
	name := slug.Make("vessel-" + cfg.Name)

	if !shutUp {
		// Ask if we can delete things
		canDeleteThings := promptui.Prompt{
			Label:     "This will permanently delete the dev environment, are you sure?",
			IsConfirm: true,
		}

		_, err = canDeleteThings.Run()

		if err != nil {
			os.Exit(0)
		}
	}

	// Attempt to stop any currently running session
	// Note that we ignore errors
	mutagen.StopSession(name)

	/**
	 * The Process:
	 * 1. Delete Fly App (which deletes machines, etc)
	 * 2. vessel.yml
	 * 3. ~/.vessel/envs/<app-name>
	 * 4. Warn about ~/.ssh/config entries (TODO: Can we safely delete from that file?)
	 */

	stopFlyctl := func() error {
		return nil
		// Does nothing, but we want it to exist, so we can call it later
		// even if we fly.ShouldStartFlyMachineApiProxy() == false
	}

	// Delete the VM if the -f / --files-only flag is not used
	// This lets you delete the VM from within Fly's and then cleanup Vessel-generated files
	if !localFiles {
		// Ensure we can connect to Fly's API
		if fly.ShouldStartFlyMachineApiProxy() {
			flyctl, err := fly.FindFlyctlCommandPath()

			if err != nil {
				logger.GetLogger().Error("command", "init", "msg", "could not find flyctl command", "error", err)
				PrintIfVerbose(Verbose, err, "You need flyctl installed to make API calls to Fly.io")

				os.Exit(1)
			}

			// Create this var, allowing us to use = instead of := in assignment below it
			// which ensures we are actually re-assigning the stopFlyctl variable
			var proxyErr error
			stopFlyctl, proxyErr = fly.StartMachineProxy(flyctl)
			time.Sleep(time.Second * 2) // Give the proxy time to boot up

			if proxyErr != nil {
				logger.GetLogger().Error("command", "init", "msg", "could not run `flyctl machine api-proxy` command", "error", err)
				PrintIfVerbose(Verbose, err, "Could not make API calls to Fly.io via api-proxy")

				os.Exit(1)
			}

			defer stopFlyctl()
		}

		err = fly.DeleteApp(auth.Token, cfg.Name)

		if err != nil {
			logger.GetLogger().Error("command", "destroy", "msg", "could not destroy Fly app", "error", err)
			PrintIfVerbose(Verbose, err, "could not destroy Fly app")
			stopFlyctl()

			os.Exit(1)
		}
	}

	err = os.Remove(ConfigPath)

	if err != nil {
		logger.GetLogger().Error("command", "destroy", "msg", "could not remove vessel project config file", "error", err)
		PrintIfVerbose(Verbose, err, "could not delete project's vessel.yml file")
		stopFlyctl()

		os.Exit(1)
	}

	appEnvDir, err := util.GetAppEnvDir(cfg.Name)

	if err != nil {
		logger.GetLogger().Error("command", "destroy", "msg", "could not find dev environment files", "error", err)
		PrintIfVerbose(Verbose, err, "could not find dev environment files in ~/.vessel/envs")
		stopFlyctl()

		os.Exit(1)
	}

	err = os.RemoveAll(appEnvDir)

	if err != nil {
		logger.GetLogger().Error("command", "destroy", "msg", "could not delete dev environment files", "error", err)
		PrintIfVerbose(Verbose, err, "could not delete dev environment files in ~/.vessel/envs")
		stopFlyctl()

		os.Exit(1)
	}

	fmt.Println("\033[1;32m\xE2\x9C\x94\033[0m Dev environment deleted")
	fmt.Printf("\033[0;33mNote:\033[0m There likely is still an entry for `vessel-%s` in your ~/.ssh/config file\n", name)
}
