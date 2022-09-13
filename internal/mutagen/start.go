package mutagen

import (
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/config"
	"strings"
)

func StartSession(name, localDir string, cfg *config.EnvironmentConfig) error {
	_, err := Sync(name, cfg.Remote.Alias, localDir, cfg.Remote.RemotePath)

	if err != nil {
		return fmt.Errorf("error starting syncing: %w", err)
	}

	// Forward multiple ports
	for k, p := range cfg.Forwarding {
		ports := strings.Split(p, ":")

		if len(ports) != 2 {
			return fmt.Errorf("invalid forwarding configuration found in vessel.yml file: %s", p)
		}

		_, err = Forward(fmt.Sprintf("%s-%d", name, k), "tcp:127.0.0.1:"+ports[0], cfg.Remote.Alias, "tcp:127.0.0.1:"+ports[1])

		if err != nil {
			return fmt.Errorf("error forwarding port config `%s`: %w", p, err)
		}
	}

	return nil
}
