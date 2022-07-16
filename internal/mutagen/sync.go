package mutagen

import (
	"encoding/json"
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/logger"
	"github.com/vessel-app/vessel-cli/internal/process"
	"os/exec"
	"strings"
)

// Sync uses Mutagen to start a sync with the given SSH and file path information
// TODO: We assume ssh alias defined in ~/.ssh/config is the only way to go
func Sync(name, alias, local_path, remote_path string) (string, error) {
	exe, err := exec.LookPath(process.ExecutableName("mutagen"))

	if err != nil {
		return "", fmt.Errorf("unable to determine mutagen path: %w", err)
	}

	// If a session of the same name exists, return that session identifier
	sessions, err := findSyncSessions(exe)

	if err != nil {
		return "", fmt.Errorf("could not list sync sessions: %w", err)
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
			"sync", "create",
			"--ignore-vcs",
			"-i", "node_modules",
			"-i", "vendor",
			"--name", name,
			"--sync-mode", "two-way-resolved",
			local_path, alias + ":" + remote_path,
		},
	}

	logger.GetLogger().Debug("sync_command", proc.String())

	output, err := proc.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("mutagen sync error: %s \n %s", exitError.String(), output)
		} else {
			return "", fmt.Errorf("unable start mutagen sync: %w", err)
		}
	}

	return string(output), nil
}

func StopSync(name string) error {
	exe, err := exec.LookPath(process.ExecutableName("mutagen"))

	if err != nil {
		return fmt.Errorf("unable to determine mutagen path: %w", err)
	}

	sessions, err := findSyncSessions(exe)

	if err != nil {
		return fmt.Errorf("could not list sync sessions: %w", err)
	}

	// It's theoretically possible to have multiple sessions of the same name
	// open, so we'll close all that match this app name
	for _, session := range sessions {
		if session.Name == name {
			stopSync := &exec.Cmd{
				Path: exe,
				Args: []string{
					exe,
					"sync", "terminate",
					session.Name,
				},
			}

			logger.GetLogger().Debug("stop_sync_command", stopSync.String())

			// Ignore errors
			stopSync.Run()
		}
	}

	return nil
}

func findSyncSessions(exe string) ([]SyncSession, error) {
	list := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe,
			"sync", "list",
			"--template", `"{{ json . }}"`,
		},
	}

	logger.GetLogger().Debug("list_syncs_command", list.String())

	output, err := list.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("mutagen list syncs error: %s \n %s", exitError.String(), output)
		} else {
			return nil, fmt.Errorf("unable to list mutagen sync sessions: %w", err)
		}
	}

	sessions := make([]SyncSession, 0)
	err = json.Unmarshal([]byte(strings.Trim(string(output), "\n\"")), &sessions)

	if err != nil {
		return nil, fmt.Errorf("could not parse mutagen output: %w", err)
	}

	return sessions, nil
}
