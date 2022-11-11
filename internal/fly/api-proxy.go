package fly

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
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

	return !isProxyRunning(false)
}

// FindFlyctlCommandPath determines if Flyctl is
// installed within the user's PATH
func FindFlyctlCommandPath() (string, error) {
	path, err := exec.LookPath("flyctl")

	if err != nil {
		return "", fmt.Errorf("could not find flyctl in PATH: %w", err)
	}

	return path, nil
}

// StartMachineProxy starts the `fly machine api-proxy` command
// and returns a function that can be used to stop it
func StartMachineProxy(exe string) (func() error, error) {
	cmd := &exec.Cmd{
		Path: exe,
		Args: []string{
			exe,
			"machine",
			"api-proxy",
		},
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve stderr for machine api-proxy: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start machine api-proxy: %w", err)
	}

	if !isProxyRunning(true) {
		proxyStderrRaw, err := io.ReadAll(stderr)
		if err != nil {
			return nil, errors.New("could not start machine api-proxy")
		}
		return nil, fmt.Errorf("could not start machine api-proxy: %s", strings.TrimSpace(string(proxyStderrRaw)))
	}

	return func() error {
		if runtime.GOOS == "windows" {
			return cmd.Process.Kill()
		}
		return cmd.Process.Signal(os.Interrupt)
	}, nil
}

// isProxyRunning determines if the fly machine API proxy is running by
// attempting to open a TCP connection to the proxy port.
func isProxyRunning(waitForStartup bool) bool {
	connectFn := func() bool {
		c, err := net.Dial("tcp", "127.0.0.1:4280")
		if err != nil {
			return false
		}
		defer c.Close()
		return true
	}
	ok := connectFn()
	if ok || !waitForStartup {
		return ok
	}
	for i := 0; i < 10; i++ {
		time.Sleep(200 * time.Millisecond)
		if connectFn() {
			return true
		}
	}

	return false
}
