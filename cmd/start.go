package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/mutagen"
	"os"
	"path/filepath"
	"strings"
)

var ConfigPath string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Sync a remote dev session",
	Long:  `Sync a forwarding and sync session based on you vessel.yml file.`,
	Run:   runStartCommand,
}

func init() {
	startCmd.Flags().StringVarP(&ConfigPath, "config-file", "c", "vessel.yml", "Configuration file to read from")
}

// runStartCommand starts Mutagen sync and forwarding sessions based on
// configuration in the vessel.yml file. It allows you to start developing remotely!
func runStartCommand(cmd *cobra.Command, args []string) {
	cfg, err := config.Retrieve(ConfigPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// TODO: SSH with Mutagen may need to edit ~/.ssh/config to add an alias for it to work reliably.
	//    This is a bit of an issue on Windows
	//    FOR NOW WE ASSUME `alias: foo` IS PROVIDED IN THE YAML FILE

	// Generate mutagen session name
	dir, err := os.Getwd()

	if err != nil {
		dir = "my-app" // todo: Something random?
	}

	// todo: We assume local path is "."
	_, err = mutagen.Sync("vessel-"+filepath.Base(dir), cfg.Remote.Alias, ".", cfg.Remote.RemotePath)

	if err != nil {
		fmt.Printf("error starting syncing session: %v", err)
		os.Exit(1)
	}

	// TODO: Forward multiple ports
	ports := strings.Split(cfg.Forwarding[0], ":")

	if len(ports) != 2 {
		fmt.Println("Port forwarding configuration must define both ports to forward separated by a colon, e.g. 8000:8000")
		os.Exit(1)
	}

	_, err = mutagen.Forward("vessel-"+filepath.Base(dir), "tcp:127.0.0.1:"+ports[0], cfg.Remote.Alias, "tcp:127.0.0.1:"+ports[1])

	if err != nil {
		fmt.Printf("error starting forwarding session: %v", err)
		os.Exit(1)
	}

	fmt.Println("Started development session")
}
