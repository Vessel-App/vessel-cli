package mutagen

import (
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/process"
	"os/exec"
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
