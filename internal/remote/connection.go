package remote

import (
	"context"
	"fmt"
	"github.com/vessel-app/vessel-cli/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Connection struct {
	config *config.RemoteConfig
}

func NewConnection(cfg *config.RemoteConfig) *Connection {
	return &Connection{
		config: cfg,
	}
}

func (c *Connection) clientConfig() (*ssh.ClientConfig, error) {
	var sshKey string
	if strings.HasPrefix(c.config.IdentityFile, "~/") {
		home, err := os.UserHomeDir()

		if err != nil {
			return nil, fmt.Errorf("cannot find home directory in ssh key search: %w", err)
		}

		sshKey = home + "/" + c.config.IdentityFile[2:]
	} else {
		sshKey = c.config.IdentityFile
	}

	key, err := ioutil.ReadFile(sshKey)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	return &ssh.ClientConfig{
		User: c.config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

func (c *Connection) Cmd(cmd string) error {
	config, err := c.clientConfig()

	if err != nil {
		return fmt.Errorf("could not create ssh client config: %w", err)
	}

	hostSocket := fmt.Sprintf("%s:%d", c.config.Hostname, c.config.Port)
	conn, err := ssh.Dial("tcp", hostSocket, config)
	if err != nil {
		return fmt.Errorf("cannot connect %v: %w", hostSocket, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %w", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	cmd = fmt.Sprintf("cd %s && %s", c.config.RemotePath, cmd)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}

func (c *Connection) SSH(ctx context.Context) error {
	config, err := c.clientConfig()

	if err != nil {
		return fmt.Errorf("could not create ssh client config: %w", err)
	}

	hostSocket := fmt.Sprintf("%s:%d", c.config.Hostname, c.config.Port)
	conn, err := ssh.Dial("tcp", hostSocket, config)
	if err != nil {
		return fmt.Errorf("cannot connect to '%v': %w", hostSocket, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %w", err)
	}
	defer session.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("terminal make raw error: %w", err)
	}
	defer terminal.Restore(fd, state)

	w, h, err := terminal.GetSize(fd)
	if err != nil {
		return fmt.Errorf("terminal get size error: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}
	if err := session.RequestPty(term, h, w, modes); err != nil {
		return fmt.Errorf("session xterm error: %w", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return fmt.Errorf("session shell error: %w", err)
	}

	if err := session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}
		return fmt.Errorf("ssh error: %w", err)
	}
	return nil
}
