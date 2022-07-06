package mutagen

import (
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/process"
	"os/exec"
)

// Sync uses Mutagen to start a sync with the given SSH and file path information
// TODO: We assume ssh alias defined in ~/.ssh/config is the only way to go
func Sync(name, alias, local_path, remote_path string) (string, error) {
	exe, err := exec.LookPath(process.ExecutableName("mutagen"))

	if err != nil {
		return "", fmt.Errorf("unable to determine mutagen path: %v", err)
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

	fmt.Printf("About to run: '%s'\n", proc.String())

	output, err := proc.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("mutagen sync error: %s \n %s", exitError.String(), output)
		} else {
			return "", fmt.Errorf("unable start mutagen sync: %v", err)
		}
	}

	return string(output), nil
}
