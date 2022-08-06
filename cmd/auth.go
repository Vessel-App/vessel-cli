package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/util"
	"github.com/vessel-app/vessel-cli/internal/vessel"
	"io/ioutil"
	"os"
	"path/filepath"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Add a authorization token",
	Long:  `Add a authorization token to generate a dev server against your account at https://vessel.app`,
	Run:   runAuthCommand,
}

var AuthToken string

func init() {
	authCmd.Flags().StringVarP(&AuthToken, "token", "t", "", "Auth token generated at https://vessel.app/user/api-tokens")
}

func runAuthCommand(cmd *cobra.Command, args []string) {
	vesselDir, err := util.MakeStorageDir()

	if err != nil {
		logger.GetLogger().Error("command", "auth", "message", "could not save auth token", "error", err)
		PrintIfVerbose(Verbose, err, "could not set auth token")

		os.Exit(1)
	}

	if len(AuthToken) == 0 {
		// Get access_token from ~/.fly/config.yml
		flycfg, err := config.RetrieveFlyConfig()

		if err != nil {
			logger.GetLogger().Error("command", "auth", "message", "could not find fly auth token", "error", err)
			PrintIfVerbose(Verbose, err, "could not find fly auth token")

			os.Exit(1)
		}

		AuthToken = flycfg.Token
	}

	// TODO: Get user & orgs
	user, err := vessel.GetUser(AuthToken)

	if err != nil {
		logger.GetLogger().Error("command", "auth", "message", "could not get user from token", "error", err)
		PrintIfVerbose(Verbose, err, "could not find a Fly user from that token")

		os.Exit(1)
	}

	var SelectedOrg vessel.Team
	if len(user.Teams) > 1 {
		selectOrg := promptui.Select{
			Label: "Which organization should we use?",
			Items: user.Teams,
			Templates: &promptui.SelectTemplates{
				Active:   fmt.Sprintf("%s {{ .Name | underline }}", promptui.IconSelect),
				Inactive: "  {{ .Name }}",
				Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Name | faint }}`, promptui.IconGood),
			},
		}

		idx, _, err := selectOrg.Run()

		if err != nil {
			// User likely bailed out
			os.Exit(1)
		}

		SelectedOrg = user.Teams[idx]
	} else {
		SelectedOrg = user.Teams[0]
	}

	yaml := fmt.Sprintf(`access_token: %s
org: %s
`, AuthToken, SelectedOrg)

	configPath := filepath.ToSlash(vesselDir + "/config.yml")
	if err = ioutil.WriteFile(configPath, []byte(yaml), 0755); err != nil {
		logger.GetLogger().Error("command", "auth", "msg", "could not write vessel config file", "error", err)
		PrintIfVerbose(Verbose, err, "could not set auth token")

		os.Exit(1)
	}

	fmt.Println("You're authenticated! Head into an application, and run `vessel init`")
}
