package config

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestCanUnmarshallValidYaml(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	got, err := RetrieveProjectConfig(cwd + "/../../testdata/vessel.yml")

	if err != nil {
		t.Errorf("could not unmarshall configuration file: %v", err)
	}

	numPortsForwarded := len(got.Forwarding)
	if numPortsForwarded != 1 {
		t.Errorf("should only see one port to forward, got %d", numPortsForwarded)
	}

	if got.Forwarding[0] != "8000:8000" {
		t.Errorf("remote port forwarding '%s' not retrieved correctly", got.Forwarding[0])
	}

	if got.Remote.Hostname != "host.example.org" {
		t.Errorf("remote hostname '%s' not retrieved correctly", got.Remote.Hostname)
	}

	if got.Remote.User != "myuser" {
		t.Errorf("remote user '%s' not retrieved correctly", got.Remote.User)
	}

	if got.Remote.IdentityFile != "~/.ssh/id_rsa" {
		t.Errorf("identity file '%s' not retrieved correctly", got.Remote.IdentityFile)
	}

	if got.Remote.Port != 22 {
		t.Errorf("port '%d' not retrieved correctly", got.Remote.Port)
	}

	if got.Remote.RemotePath != "/home/myuser/myapp" {
		t.Errorf("remote path '%s' not retrieved correctly", got.Remote.RemotePath)
	}

	valid, err := got.Valid()
	if !valid {
		t.Errorf("yaml configuration file was not valid: %s", err)
	}
}

func TestInValidYamlIsFoundToNotBeNotValid(t *testing.T) {
	cwd, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	_, invalidErr := RetrieveProjectConfig(cwd + "/../../testdata/invalid.yml")

	if invalidErr == nil {
		t.Error("yaml validity should have returned an error")
	}

	if !strings.Contains(invalidErr.Error(), "user") {
		t.Errorf("invalid yaml error should be specific to missing username: %v", err)
	}
}
