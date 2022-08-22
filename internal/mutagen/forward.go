package mutagen

import (
	"encoding/json"
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"os/exec"
	"strings"
)

// Forward uses Mutagen to start a forward with the given ports
// TODO: We assume ssh alias defined in ~/.ssh/config is the only way to go
func Forward(name, localSocket, alias, remoteSocket string) (string, error) {
	exe, err := GetMutagenPath()

	if err != nil {
		return "", fmt.Errorf("unable to determine mutagen path: %v", err)
	}

	// If a session of the same name exists, return that session identifier
	sessions, err := findForwardSessions(exe)

	if err != nil {
		return "", fmt.Errorf("could not list forward sessions: %w", err)
	}

	for _, session := range sessions {
		if session.Name == name {
			return session.Identifier, nil
		}
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
			return "", fmt.Errorf("mutagen forward error: %s \n %s", exitError.Error(), output)
		} else {
			return "", fmt.Errorf("unable start mutagen forward: %w", err)
		}
	}

	return string(output), nil
}

func StopForward(name string) error {
	exe, err := GetMutagenPath()

	if err != nil {
		return fmt.Errorf("unable to determine mutagen path: %w", err)
	}

	sessions, err := findForwardSessions(exe)

	if err != nil {
		return fmt.Errorf("could not list forward sessions: %w", err)
	}

	// It's theoretically possible to have multiple sessions of the same name
	// open, so we'll close all that match this app name
	for _, session := range sessions {
		if session.Name == name {
			stopForward := &exec.Cmd{
				Path: exe,
				Args: []string{
					exe,
					"forward", "terminate",
					session.Name,
				},
			}

			logger.GetLogger().Debug("stop_sync_command", stopForward.String())

			// Ignore errors
			stopForward.Run()
		}
	}

	return nil
}

// findSessions will find any forward/sync
func findForwardSessions(exe string) ([]ForwardSession, error) {
	list := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe,
			"forward", "list",
			"--template", `"{{ json . }}"`,
		},
	}

	logger.GetLogger().Debug("list_forwards_command", list.String())

	output, err := list.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("mutagen list forwards error: %s \n %s", exitError.Error(), output)
		} else {
			return nil, fmt.Errorf("unable to list mutagen forward sessions: %w", err)
		}
	}

	sessions := make([]ForwardSession, 0)
	err = json.Unmarshal([]byte(strings.Trim(string(output), "\n\"")), &sessions)

	if err != nil {
		return nil, fmt.Errorf("could not parse mutagen output: %w", err)
	}

	return sessions, nil
}
