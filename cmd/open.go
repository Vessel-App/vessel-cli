package cmd

import (
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"os"
	"strings"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open the browser for the current dev environment",
	Long:  `Open the browser for the current dev environment`,
	Run:   runOpenCommand,
}

func runOpenCommand(cmd *cobra.Command, args []string) {
	cfg, err := config.Retrieve(ConfigPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Best attempt at guessing the port.
	// This will be an annoying bug report some day.
	if len(cfg.Forwarding) > 0 {
		portCombo := strings.Split(cfg.Forwarding[0], ":")

		err := open.Run("http://localhost:" + portCombo[0])

		if err != nil {
			logger.GetLogger().Error("command", "open", "msg", "could not run open command", "error", err)
			PrintIfVerbose(Verbose, err, "error opening a browser")

			fmt.Printf("error opening browser: %v", err)
			os.Exit(1)
		}
	} else {
		logger.GetLogger().Error("command", "open", "msg", "No port forwarding defined, cannot open your browser", "error", err)
		PrintIfVerbose(Verbose, err, "error opening a browser")

		os.Exit(1)
	}
}
