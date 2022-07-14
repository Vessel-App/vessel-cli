package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Reference: https://zetcode.com/golang/yaml/
func Retrieve(path string) (*EnvironmentConfig, error) {
	file, err := ioutil.ReadFile(path)

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
