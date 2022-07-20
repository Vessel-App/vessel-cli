package cmd

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/kevinburke/ssh_config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/util"
	"io/ioutil"
	"os"
	"path/filepath"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure an environment.",
	Long:  `Configure and create a remote dev environment.`,
	Run:   runInitCommand,
}

var AppName string

func init() {
	initCmd.Flags().StringVarP(&AppName, "name", "n", "", "Set the environment name")
}

func runInitCommand(cmd *cobra.Command, args []string) {
	auth, err := config.RetrieveVesselConfig()

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not get vessel config", "error", err)
		PrintIfVerbose(Verbose, err, "init error, please make sure to run `vessel auth` first")

		os.Exit(1)
	}

	// Get/generate application name
	dir, err := os.Getwd()

	if err != nil {
		dir = "my-app"
	}

	askAppName := promptui.Prompt{
		Label:   "App Name",
		Default: slug.Make(filepath.Base(dir)),
	}

	appName, err := askAppName.Run()

	if err != nil {
		// No logging, user likely just bailed out
		os.Exit(1)
	}

	appName = slug.Make(appName)

	// TODO: If this app exists, ask if we should over-write stuff

	// Create ~/.vessel/<app-name>
	vesselAppDir, err := util.MakeAppDir(appName)

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not create vessel storage directory", "error", err, "dir", vesselAppDir)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	// Generate and store SSH keys
	keys, err := util.GenerateSSHKey()

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not generate SSH keys", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	privateKeyPath := filepath.FromSlash(vesselAppDir + "/id_ed25519")
	if err = ioutil.WriteFile(privateKeyPath, keys.Private, 0600); err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not store generated SSH private key", "error", err, "file", privateKeyPath)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	publicKeyPath := filepath.FromSlash(vesselAppDir + "/id_ed25519.pub")
	if err = ioutil.WriteFile(publicKeyPath, keys.Public, 0644); err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not store generated SSH public key", "error", err, "file", publicKeyPath)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	// TODO: Send public key back to Vessel API for server creation
	//       (And wait for the server to come alive?)
	fmt.Printf("HERE we should call the API with token %s and have it create a machine for team %s.", auth.Token, auth.TeamGuid)

	// Ask if we can add to ~/.ssh/config (if alias is not present)
	canAddSSHAlias := promptui.Prompt{
		Label:     "Can we add a Host entry to your ~/.ssh/config file ('N' will output to terminal instead)",
		IsConfirm: true,
	}

	sshConfig := fmt.Sprintf(`
Host vessel-%s
    HostName TODO
    User vessel
    IdentityFile %s
    IdentitiesOnly yes
`, appName, privateKeyPath)

	_, err = canAddSSHAlias.Run()

	if err != nil {
		fmt.Println("Here is what we would have added to ~/.ssh/config:")
		fmt.Println(sshConfig)
	} else {
		// Only write if host doesn't exist yet
		hostAlreadyExists := ssh_config.Get(appName, "HostName")
		if len(hostAlreadyExists) == 0 {
			if err = util.WriteToSshConfig(sshConfig); err != nil {
				logger.GetLogger().Error("command", "init", "msg", "could not write to SSH config to ~/.ssh/config", "error", err)
				PrintIfVerbose(Verbose, err, "error initializing app")

				os.Exit(1)
			}
		} else {
			fmt.Printf("Warning: ~/.ssh/config file already contained Host %s", appName)
		}
	}

	yaml := fmt.Sprintf(`name: %s

remote:
  hostname: "TODO: Set this"
  user: vessel
  identityfile: %s
  port: 2222
  path: /home/vessel/app
  alias: vessel-%s

forwarding:
  - 8000:80
`, appName, privateKeyPath, appName)

	if err = ioutil.WriteFile("vessel.yml", []byte(yaml), 0755); err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not write yaml file to current directory", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	fmt.Println("You're good to go! Run `vessel start` to begin developing!")
}
