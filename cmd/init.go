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
	"github.com/vessel-app/vessel-cli/internal/remote"
	"github.com/vessel-app/vessel-cli/internal/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a dev environment.",
	Long: `Configure and create a remote dev environment.
Defaults to assigning an IPv6 address. Use the -4 flag to use IPv4 instead.`,
	Run: runInitCommand,
}

var AppName string
var UseIpv4 bool

func init() {
	initCmd.Flags().StringVarP(&AppName, "name", "n", "", "Set the environment name")
	initCmd.Flags().BoolVarP(&UseIpv4, "ipv4", "4", false, "Allocate an IPv4 instead of IPv6 to the environment")
}

// runInitCommand will guide users through setting up a new development environment.
// It performs the following actions:
//  1. Retrieves vessel configuration
//  2. Starts `fly machine api-proxy` if needed
//  3. Helps create an environment name
//  4. Prompts for dev env type (PHP, etc)
//  5. Generates env files (SSH keys, etc)
//  6. Gets user's nearest region
//  7. Creates the dev environment
//  8. Generates project and SSH configuration
//  9. Downloads Mutagen (if needed)
//  10. Waits for dev env to be available
func runInitCommand(cmd *cobra.Command, args []string) {
	auth, err := config.RetrieveVesselConfig()

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not get vessel config", "error", err)
		PrintIfVerbose(Verbose, err, "init error, please make sure to run `vessel auth` first")

		os.Exit(1)
	}

	stopFlyctl := func() {
		// Does nothing, but we want it to exist so we can call it later
		// even if we fly.ShouldStartFlyMachineApiProxy() == false
	}

	// Ensure we can connect to Fly's API
	if fly.ShouldStartFlyMachineApiProxy() {
		flyctl, err := fly.FindFlyctlCommandPath()

		if err != nil {
			logger.GetLogger().Error("command", "init", "msg", "could not find flyctl command", "error", err)
			PrintIfVerbose(Verbose, err, "You need flyctl installed to make API calls to Fly.io")

			os.Exit(1)
		}

		stopFlyctl, err := fly.StartMachineProxy(flyctl)

		if err != nil {
			logger.GetLogger().Error("command", "init", "msg", "could not run `flyctl machine api-proxy` command", "error", err)
			PrintIfVerbose(Verbose, err, "Could not make API calls to Fly.io via api-proxy")

			os.Exit(1)
		}

		defer stopFlyctl()
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
		stopFlyctl()
		os.Exit(1)
	}

	appName = slug.Make(appName)

	// Get image to use (development environment type)
	bundledTypes := []string{"vesselapp/php:8.1", "vesselapp/php:8.0"}
	typeIndex := -1
	var envDockerImage string

	for typeIndex < 0 {
		typePrompt := promptui.SelectWithAdd{
			Label:    "What Docker base image should we use?",
			Items:    bundledTypes,
			AddLabel: "Other",
		}

		typeIndex, envDockerImage, err = typePrompt.Run()

		if typeIndex == -1 {
			bundledTypes = append(bundledTypes, envDockerImage)
		}
	}

	if err != nil {
		logger.GetLogger().Debug("cmd", "init", "msg", "prompt ui failure selecting Docker image", "error", err)
		stopFlyctl()
		os.Exit(1)
	}

	// Create ~/.vessel/<app-name>
	vesselAppDir, err := util.MakeAppDir(appName)

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not create vessel storage directory", "error", err, "dir", vesselAppDir)
		PrintIfVerbose(Verbose, err, "error initializing app")
		stopFlyctl()
		os.Exit(1)
	}

	// Generate and store SSH keys
	keys, err := util.GenerateSSHKey()

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not generate SSH keys", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")
		stopFlyctl()
		os.Exit(1)
	}

	privateKeyPath := filepath.FromSlash(vesselAppDir + "/id_ed25519")
	if err = ioutil.WriteFile(privateKeyPath, keys.Private, 0600); err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not store generated SSH private key", "error", err, "file", privateKeyPath)
		PrintIfVerbose(Verbose, err, "error initializing app")
		stopFlyctl()
		os.Exit(1)
	}

	publicKeyPath := filepath.FromSlash(vesselAppDir + "/id_ed25519.pub")
	if err = ioutil.WriteFile(publicKeyPath, keys.Public, 0644); err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not store generated SSH public key", "error", err, "file", publicKeyPath)
		PrintIfVerbose(Verbose, err, "error initializing app")
		stopFlyctl()
		os.Exit(1)
	}

	// Get user's nearest Fly region
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
			stopFlyctl()
			os.Exit(1)
		}

		nearestRegionCode = fly.Regions[idx].Code
	} else {
		nearestRegionCode = region.NearestRegion.Code
	}

	// Create dev environment
	env, err := environments.CreateEnvironment(auth.Token, appName, envDockerImage, auth.Org, nearestRegionCode, string(keys.Public), !UseIpv4)

	if err != nil {
		logger.GetLogger().Debug("cmd", "init", "msg", "could not create dev environment", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")
		stopFlyctl()
		os.Exit(1)
	}

	w := wow.New(os.Stdout, spin.Get(spin.Dots), " Environment registered, waiting for it to start")
	w.Start()

	err = fly.WaitForMachine(auth.Token, appName, env.FlyMachine)
	if err != nil {
		logger.GetLogger().Debug("cmd", "init", "msg", "could not get dev environment status", "error", err)
		PrintIfVerbose(Verbose, err, "error creating dev environment")
		stopFlyctl()
		os.Exit(1)
	}

	w.PersistWith(spin.Spinner{Frames: []string{"\033[1;32m\xE2\x9C\x94\033[0m"}}, " Environment ready!")

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
    UserKnownHostsFile /dev/null
    StrictHostKeyChecking no
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
				stopFlyctl()
				os.Exit(1)
			}
		} else {
			fmt.Printf("Warning: ~/.ssh/config file already contained Host %s", appName)
		}
	}

	// Generate project configuration file
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
		stopFlyctl()
		os.Exit(1)
	}

	// Ensure Mutagen is installed
	w2 := wow.New(os.Stdout, spin.Get(spin.Dots), " Configuring Mutagen")
	w2.Start()

	_, err = util.MakeBinDir()

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not create ~/.vessel/bin directory", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")
		stopFlyctl()
		os.Exit(1)
	}

	err = mutagen.InstallMutagen()

	if err != nil {
		logger.GetLogger().Error("command", "init", "msg", "could not install mutagen to ~/.vessel/bin/mutagen", "error", err)
		PrintIfVerbose(Verbose, err, "error initializing app")
		stopFlyctl()
		os.Exit(1)
	}

	w2.PersistWith(spin.Spinner{Frames: []string{"\033[1;32m\xE2\x9C\x94\033[0m"}}, " Mutagen is ready")

	// Wait for SSH to become available
	w3 := wow.New(os.Stdout, spin.Get(spin.Dots), " Waiting for environment to become reachable")
	w3.Start()

	cfg, err := config.RetrieveProjectConfig(ConfigPath)

	if err != nil {
		logger.GetLogger().Error("command", "init", "error", err)
		PrintIfVerbose(Verbose, err, "could not read configuration")
		stopFlyctl()
		os.Exit(1)
	}

	connection := remote.NewConnection(&cfg.Remote)

	// Wait for up to ~30 seconds for SSH to become available
	// (15 attempts, attempted every 2 seconds)
	attempts := 0
	success := false
	for attempts <= 15 {
		if err := connection.TestConnection(); err != nil {
			time.Sleep(time.Second * 2)
			attempts++
		} else {
			success = true
			break
		}
	}

	if !success {
		logger.GetLogger().Error("command", "init", "error", err)
		PrintIfVerbose(Verbose, err, "could not connect to dev environment")
		stopFlyctl()
		os.Exit(1)
	}

	w3.PersistWith(spin.Spinner{Frames: []string{"\033[1;32m\xE2\x9C\x94\033[0m"}}, " Environment is reachable")

	fmt.Println("You're good to go! Run `vessel start` to begin developing!")
}
