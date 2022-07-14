package util

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

// MakeStorageDir creates a ~/.vessel directory
func MakeStorageDir() (string, error) {
	home, err := homedir.Dir()

	if err != nil {
		return "", fmt.Errorf("could not find home dir: %w", err)
	}

	vesselPath := filepath.FromSlash(home + "/.vessel")

	// Check if the ~/.vessel directory exists already
	// We assume a home directory always exists
	stat, err := os.Stat(vesselPath)

	if err == nil && stat.IsDir() {
		// path exists already
		return vesselPath, nil
	} else if err == nil && !stat.IsDir() {
		// path exists but is not a directory
		return "", fmt.Errorf("could not create directory ~/.vessel as a file with that name already exists")
	} else if os.IsNotExist(err) {
		// path does not exist, create it
		if err := os.Mkdir(vesselPath, 0750); err != nil {
			return "", fmt.Errorf("could not create vessel directory: %w", err)
		}

		return vesselPath, nil
	} else if err != nil {
		return "", fmt.Errorf("stat error: %w", err)
	}

	return vesselPath, nil
}

// MakeAppDir creates a ~/.vessel/<app-name> directory
func MakeAppDir(appName string) (string, error) {
	vesselPath, err := MakeStorageDir()

	if err != nil {
		return "", fmt.Errorf("could not create vessel storage dir: %w", err)
	}

	vesselAppPath := filepath.FromSlash(vesselPath + "/" + appName)

	stat, err := os.Stat(vesselAppPath)

	if err == nil && stat.IsDir() {
		// path exists already
		return vesselAppPath, nil
	} else if err == nil && !stat.IsDir() {
		// path exists but is not a directory
		return "", fmt.Errorf("could not create vessel app directory as a file with that name already exists")
	} else if os.IsNotExist(err) {
		// path does not exist, create it
		if err := os.Mkdir(vesselAppPath, 0750); err != nil {
			return "", fmt.Errorf("could not create vessel app directory: %w", err)
		}

		return vesselAppPath, nil
	} else if err != nil {
		return "", fmt.Errorf("stat error: %w", err)
	}

	return vesselAppPath, nil
}
