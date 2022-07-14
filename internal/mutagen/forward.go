package mutagen

import (
	"encoding/json"
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/process"
	"os/exec"
	"strings"
)

// Forward uses Mutagen to start a forward with the given ports
// TODO: We assume ssh alias defined in ~/.ssh/config is the only way to go
func Forward(name, localSocket, alias, remoteSocket string) (string, error) {
	exe, err := exec.LookPath(process.ExecutableName("mutagen"))

	if err != nil {
		return "", fmt.Errorf("unable to determine mutagen path: %v", err)
	}

	proc := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe,
			"forward", "create",
			"--name", name,
			localSocket, alias + ":" + remoteSocket,
		},
	}

	logger.GetLogger().Debug("forward_command", proc.String())

	output, err := proc.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("mutagen forward error: %s \n %s", exitError.String(), output)
		} else {
			return "", fmt.Errorf("unable start mutagen forward: %w", err)
		}
	}

	return string(output), nil
}

func StopForward(name string) error {
	exe, err := exec.LookPath(process.ExecutableName("mutagen"))

	if err != nil {
		return fmt.Errorf("unable to determine mutagen path: %w", err)
	}

	list := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe,
			"forward", "list",
			"--template", `"{{ json . }}"`,
		},
	}

	logger.GetLogger().Debug("stop_sync_command", list.String())

	output, err := list.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("mutagen sync error: %s \n %s", exitError.String(), output)
		} else {
			return fmt.Errorf("unable to stop mutagen sync: %w", err)
		}
	}

	sessions := make([]ForwardSession, 0)
	err = json.Unmarshal([]byte(strings.Trim(string(output), "\n\"")), &sessions)

	if err != nil {
		return fmt.Errorf("could not parse mutagen output: %w", err)
	}

	for _, session := range sessions {
		if session.Name == name {
			// Run immediately, ignoring errors
			(&exec.Cmd{
				Path: exe,
				Args: []string{
					exe,
					"sync", "terminate",
					session.Name,
				},
			}).Run()
		}
	}

	return nil
}
