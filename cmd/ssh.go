package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vessel-app/vessel-cli/internal/config"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/remote"
	"os"
	"os/signal"
	"syscall"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Log into the remote dev environment",
	Long:  `Start an SSH session in the remove dev environment. Poke around!`,
	Run:   runSSHCommand,
}

// runSSHCommand starts an interactive SSH session
// with the remote development environment.
func runSSHCommand(cmd *cobra.Command, args []string) {
	cfg, err := config.Retrieve(ConfigPath)

	if err != nil {
		logger.GetLogger().Error("command", "ssh", "msg", "could not read configuration", "error", err)
		PrintIfVerbose(Verbose, err, "error running SSH")

		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := run(ctx, cfg); err != nil {
			logger.GetLogger().Error("command", "ssh", "msg", "error running SSH", "error", err)
			PrintIfVerbose(Verbose, err, "error running SSH")

			os.Exit(1)
		}
		cancel()
	}()

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
}

func run(ctx context.Context, cfg *config.EnvironmentConfig) error {
	connection := remote.NewConnection(&cfg.Remote)

	err := connection.SSH(ctx)
	if err != nil {
		return fmt.Errorf("could not start ssh session: %w", err)
	}

	return nil
}
