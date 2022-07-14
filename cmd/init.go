package cmd

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/kevinburke/ssh_config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
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
		os.Exit(1)
	}

	appName = slug.Make(appName)

	// TODO: If this app exists, ask if we should over-write stuff

	// Create ~/.vessel/<app-name>
	vesselAppDir, err := util.MakeAppDir(appName)

	if err != nil {
		fmt.Printf("could not create vessel storage directory '%s': %v", vesselAppDir, err)
		os.Exit(1)
	}

	// Generate and store SSH keys
	keys, err := util.GenerateSSHKey()

	if err != nil {
		fmt.Printf("ssh key error: %v", err)
		os.Exit(1)
	}

	privateKeyPath := filepath.FromSlash(vesselAppDir + "/id_ed25519")
	if err = ioutil.WriteFile(privateKeyPath, keys.Private, 0600); err != nil {
		fmt.Printf("could not store generated SSH private key: %v", err)
		os.Exit(1)
	}

	if err = ioutil.WriteFile(filepath.FromSlash(vesselAppDir+"/id_ed25519.pub"), keys.Public, 0644); err != nil {
		fmt.Printf("could not store generated SSH public key: %v", err)
		os.Exit(1)
	}

	// TODO: Send public key back to Vessel API for server creation

	// Ask if we can add to ~/.ssh/config (if alias is not present)
	canAddSSHAlias := promptui.Prompt{
		Label:     "We need to add a Host entry to your ~/.ssh/config file, is that okay? Using 'N' will output what we would add there instead.",
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
		fmt.Println("No problem! Here is what we would have added to ~/.ssh/config")
		fmt.Println(sshConfig)
	} else {
		// Only write if host doesn't exist yet
		hostAlreadyExists := ssh_config.Get(appName, "HostName")
		if len(hostAlreadyExists) == 0 {
			if err = util.WriteToSshConfig(sshConfig); err != nil {
				fmt.Printf("could not write to SSH config file: %v", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("SSH config already contains host %s", appName)
			os.Exit(1)
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
		fmt.Printf("could not write yaml file to current directory: %v", err)
		os.Exit(1)
	}

	fmt.Println("You're good to go! Run `vessel start` to begin developing!")
}
