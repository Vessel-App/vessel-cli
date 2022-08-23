package cmd

import (
	"fmt"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/gosimple/slug"
	"github.com/kevinburke/ssh_config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/environments"
	"github.com/vessel-app/vessel-cli/internal/fly"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/mutagen"
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
		logger.GetLogger().Debug("cmd", "init", "msg", "prompt ui failure asking app name", "error", err)
		os.Exit(1)
	}

	appName = slug.Make(appName)

	// TODO: If this app exists, ask if we should over-write ~/.vessel/<app-name> files

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

	var nearestRegionCode string
	region, err := fly.GetNearestRegion(auth.Token)

	if err != nil {
		logger.GetLogger().Debug("cmd", "init", "msg", "could not automatically find nearest region", "error", err)

		// If we can't get the nearest region, have them select a region
		selectRegion := promptui.Select{
			Label: "Which region is closest to you?",
			Items: fly.Regions,
			Templates: &promptui.SelectTemplates{
				Active:   fmt.Sprintf("%s {{ .Code | underline }}{{ `-` | underline }}{{ .Name | underline }}", promptui.IconSelect),
				Inactive: "  {{ .Code }} - {{ .Name }}",
				Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Code| faint }}{{ "-" | faint }}{{ .Name | faint }}`, promptui.IconGood),
			},
			Size: len(fly.Regions),
		}

		idx, _, err := selectRegion.Run()

		if err != nil {
			// User likely bailed out
			logger.GetLogger().Debug("cmd", "init", "msg", "prompt ui failure selecting region", "error", err)
			os.Exit(1)
		}

		nearestRegionCode = fly.Regions[idx].Code
	} else {
		nearestRegionCode = region.NearestRegion.Code
	}

	env, err := environments.CreateEnvironment(auth.Token, appName, auth.Org, nearestRegionCode, string(keys.Public))

	if err != nil {
		logger.GetLogger().Debug("cmd", "init", "msg", "could not create dev environment", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	w := wow.New(os.Stdout, spin.Get(spin.Dots), " Environment registered, waiting for it to become available")
	w.Start()

	err = fly.WaitForMachine(auth.Token, appName, env.FlyMachine)
	if err != nil {
		logger.GetLogger().Debug("cmd", "init", "msg", "could not get dev environment status", "error", err)
		PrintIfVerbose(Verbose, err, "error creating dev environment")

		os.Exit(1)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"👍"}}, " Environment ready!")

	// Ask if we can add to ~/.ssh/config (if alias is not present)
	canAddSSHAlias := promptui.Prompt{
		Label:     "Can we add a Host entry to your ~/.ssh/config file ('N' will output to terminal instead)",
		IsConfirm: true,
	}

	sshConfig := fmt.Sprintf(`
Host vessel-%s
    HostName %s
    User vessel
    IdentityFile %s
    IdentitiesOnly yes
    AddressFamily inet6
`, appName, env.FlyIp, privateKeyPath)

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
  hostname: %s
  user: vessel
  identityfile: %s
  port: 22
  path: /home/vessel/app
  alias: vessel-%s

forwarding:
  - 8000:80
`, appName, env.FlyIp, privateKeyPath, appName)

	if err = os.WriteFile("vessel.yml", []byte(yaml), 0755); err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not write yaml file to current directory", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	w2 := wow.New(os.Stdout, spin.Get(spin.Dots), " Configuring Mutagen")
	w2.Start()

	_, err = util.MakeBinDir()

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not create ~/.vessel/bin directory", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	err = mutagen.InstallMutagen()

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not install mutagen to ~/.vessel/bin/mutagen", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")

		os.Exit(1)
	}

	w2.PersistWith(spin.Spinner{Frames: []string{"👍"}}, " Mutagen installed")

	fmt.Println("You're good to go! Run `vessel start` to begin developing!")
}
