package mutagen

import (
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"os/exec"
)

// StopMutagenDaemon will try to stop a running mutagen daemon
func StopMutagenDaemon() (string, error) {
	exe, err := GetMutagenCommandPath()

	if err != nil {
		return "", fmt.Errorf("unable to determine mutagen path: %w", err)
	}

	proc := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe,
			"daemon", "stop",
		},
	}

	logger.GetLogger().Debug("stop_daemon_command", proc.String())

	output, err := proc.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("mutagen stop daemon error: %s \n %s", exitError.Error(), output)
		} else {
			return "", fmt.Errorf("unable stop mutagen daemon: %w", err)
		}
	}

	return string(output), nil
}
