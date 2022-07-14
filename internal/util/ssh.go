package util

import (
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"github.com/mikesmitty/edkey"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
)

type Keys struct {
	Public, Private []byte
}

func GenerateSSHKey() (*Keys, error) {
	// Generate a new private/public keypair for OpenSSH
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		return nil, fmt.Errorf("could not generate ssh keys: %w", err)
	}

	publicKey, err := ssh.NewPublicKey(pubKey)

	if err != nil {
		return nil, fmt.Errorf("could not generate public key: %w", err)
	}

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(privKey),
	}

	return &Keys{
		Public:  ssh.MarshalAuthorizedKey(publicKey),
		Private: pem.EncodeToMemory(pemKey),
	}, nil
}

func WriteToSshConfig(content string) error {
	home, err := homedir.Dir()

	if err != nil {
		return fmt.Errorf("could not find home dir: %w", err)
	}

	f, err := os.OpenFile(filepath.FromSlash(home+"/.ssh/config"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		return fmt.Errorf("could not open ssh config file for writing: %w", err)
	}

	defer f.Close()

	if _, err = f.WriteString(content); err != nil {
		return fmt.Errorf("could not write to ssh config file: %w", err)
	}

	return nil
}
