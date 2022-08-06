package config

import "fmt"

type AuthConfig struct {
	Token    string `yaml:"access_token"`
	TeamGuid string `yaml:"team"`
}

type EnvironmentConfig struct {
	Name       string       `yaml:"name"`
	Remote     RemoteConfig `yaml:"remote"`
	Forwarding []string     `yaml:"forwarding"`
}

type RemoteConfig struct {
	Hostname     string `yaml:"hostname"`
	User         string `yaml:"user"`
	IdentityFile string `yaml:"identityfile"`
	Port         int    `yaml:"port"`
	RemotePath   string `yaml:"path"`
	Alias        string `yaml:alias,omitempty`
}

func (c *EnvironmentConfig) Valid() (bool, error) {
	if len(c.Name) < 1 {
		return false, fmt.Errorf("no app name defined")
	}

	if len(c.Forwarding) < 1 {
		return false, fmt.Errorf("no forwarding ports are defined")
	}

	if c.Remote.Hostname == "" {
		return false, fmt.Errorf("no remote hostname (to SSH in with) defined")
	}

	if c.Remote.User == "" {
		return false, fmt.Errorf("no remote user (to SSH in with) defined")
	}

	if c.Remote.IdentityFile == "" {
		return false, fmt.Errorf("no remote indentity file (SSH key) defined")
	}

	if c.Remote.RemotePath == "" {
		return false, fmt.Errorf("no remote file path to sync with defined")
	}

	if c.Remote.Port == 0 {
		return false, fmt.Errorf("no remote port (to SSH into) defined")
	}

	return true, nil
}
