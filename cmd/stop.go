package cmd

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
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
	cfg, err := config.Retrieve(ConfigPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Generate mutagen session name
	// Get mutagen session name
	name := slug.Make("vessel-" + cfg.Name)

	// stop syncing
	errSync := mutagen.StopSync(name)

	// stop forwarding
	errForward := mutagen.StopForward(name)

	if errSync != nil || errForward != nil {
		fmt.Println("Error disconnecting from development server")
		os.Exit(1)
	}

	fmt.Println("Stopped development session")
}
