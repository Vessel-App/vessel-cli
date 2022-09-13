package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// Reference: https://zetcode.com/golang/yaml/

// RetrieveProjectConfig will find and parse a vessel.yml file for a given project
func RetrieveProjectConfig(path string) (*EnvironmentConfig, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("could not read yaml file '%s': %w", path, err)
	}

	cfg := &EnvironmentConfig{}

	unMarshallErr := yaml.Unmarshal(file, cfg)

	if unMarshallErr != nil {
		return nil, fmt.Errorf("error parsing yaml file %s: %w", path, unMarshallErr)
	}

	valid, err := cfg.Valid()

	if !valid {
		return nil, fmt.Errorf("invalid yaml configuration: %w", err)
	}

	return cfg, nil
}

func RetrieveVesselConfig() (*AuthConfig, error) {
	home, err := homedir.Dir()

	if err != nil {
		return nil, fmt.Errorf("could not retrieve vessel config: %w", err)
	}

	configPath := filepath.ToSlash(home + "/.vessel/config.yml")
	file, err := os.ReadFile(configPath)

	if err != nil {
		return nil, fmt.Errorf("could not read yaml file '%s': %w", configPath, err)
	}

	cfg := &AuthConfig{}

	err = yaml.Unmarshal(file, cfg)

	if err != nil {
		return nil, fmt.Errorf("error parsing yaml file %s: %w", configPath, err)
	}

	return cfg, nil
}

func RetrieveFlyConfig() (*FlyConfig, error) {
	home, err := homedir.Dir()

	if err != nil {
		return nil, fmt.Errorf("could not retrieve fly config: %w", err)
	}

	configPath := filepath.ToSlash(home + "/.fly/config.yml")
	file, err := os.ReadFile(configPath)

	if err != nil {
		return nil, fmt.Errorf("could not read yaml file '%s': %w", configPath, err)
	}

	cfg := &FlyConfig{}

	err = yaml.Unmarshal(file, cfg)

	if err != nil {
		return nil, fmt.Errorf("error parsing yaml file %s: %w", configPath, err)
	}

	return cfg, nil
}
