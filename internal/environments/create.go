package environments

import (
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/fly"
)

type Environment struct {
	FlyApp     string
	FlyOrg     string
	FlyIp      string
	FlyMachine string
}

func CreateEnvironment(token, appName, image, org, region, pubKey string) (*Environment, error) {
	// Create App
	app, err := fly.CreateApp(token, appName, org)

	if err != nil {
		return nil, fmt.Errorf("could not register app: %w", err)
	}

	// Run Machine (image + env var)
	machine, err := fly.RunMachine(token, appName, region, image, pubKey)

	if err != nil {
		return nil, fmt.Errorf("could not run machine: %w", err)
	}

	// Allocate IP
	ip, err := fly.AllocateIp(token, appName, true)

	if err != nil {
		return nil, fmt.Errorf("could not allocate ip: %w", err)
	}

	return &Environment{
		FlyApp:     app.AppName,
		FlyOrg:     org,
		FlyIp:      ip.IpAddress.Address,
		FlyMachine: machine.Id,
	}, nil
}
