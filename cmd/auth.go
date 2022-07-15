package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/util"
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
	authCmd.MarkFlagRequired("token")
}

func runAuthCommand(cmd *cobra.Command, args []string) {
	vesselDir, err := util.MakeStorageDir()

	if err != nil {
		logger.GetLogger().Error("command", "auth", "message", "could not save auth token", "error", err)
		PrintIfVerbose(Verbose, err, "could not set auth token")

		os.Exit(1)
	}

	yaml := fmt.Sprintf(`access_token: %s
`, AuthToken)

	configPath := filepath.ToSlash(vesselDir + "/config.yml")
	if err = ioutil.WriteFile(configPath, []byte(yaml), 0755); err != nil {
		logger.GetLogger().Error("command", "auth", "msg", "could not write vessel config file", "error", err)
		PrintIfVerbose(Verbose, err, "could not set auth token")

		os.Exit(1)
	}

	fmt.Println("You're authenticated! Head into an application, and run `vessel init`")
}
