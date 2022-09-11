package fly

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)

// ShouldStartFlyMachineApiProxy will attempt to run the `fly machine api-proxy` command
// if FLY_HOST env var is not set and no connection to 127.0.0.1:4280 can be made (user
// did not already start the machine api-proxy).
func ShouldStartFlyMachineApiProxy() bool {
	// If FLY_HOST is set, we do nothing, as we assume
	// the user is VPNed into Fly and using _api.internal:4280
	flyHost := os.Getenv("FLY_HOST")

	if len(flyHost) > 0 {
		flyApiHost = flyHost
		return false
	}

	// Test if we can connect to the proxy address
	conn, err := net.Dial("tcp", "127.0.0.1:4280")
	if err != nil {
		// If proxying is not happening unable to connect,
		// we *should* attempt to start the proxy
		return true
	}
	defer conn.Close()

	// We were able to connect, user likely already has
	// the machine api-proxy running
	return false
}

// FindFlyctlCommandPath determines if Flyctl is installed
// in the user's PATH
func FindFlyctlCommandPath() (string, error) {
	path, err := exec.LookPath("flyctl")

	if err != nil {
		return "", fmt.Errorf("could not find flyctl in PATH: %w", err)
	}

	return path, nil
}

// TODO: Need to start this in the background and stop it later
//       Run it in a goroutine and have a function to stop it (how to stop it?)
func StartMachineProxy() error {
	proc := &exec.Cmd{
		Path: "flyctl",
		Args: []string{
			"fyctl",
			"machine",
			"api-proxy",
		},
	}
	return nil
}

func StopMachineProxy() error {
	return nil
}
