package pacman

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const (
	timeout           = 10 * time.Second
	pacmanCommand     = "pacman"
	pacmanConfCommand = "pacman-conf"
	repository        = "core"
)

type CommandExecutor interface {
	Execute(name string, arg ...string) ([]byte, error)
}

type osCommandExecutor struct{}

func (r osCommandExecutor) Execute(name string, arg ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return exec.CommandContext(ctx, name, arg...).Output()
}

type Pacman struct {
	Executor CommandExecutor
}

func NewPacman() Pacman {
	return Pacman{Executor: osCommandExecutor{}}
}

func (pacman *Pacman) GetInstalledPackages() ([]string, error) {
	out, err := pacman.Executor.Execute(pacmanCommand, "-Qq")

	return strings.Split(strings.TrimSpace(string(out)), "\n"), err
}

func (pacman *Pacman) GetServer() (string, error) {
	out, err := pacman.Executor.Execute(pacmanConfCommand, "--repo", repository, "Server")
	if err != nil {
		return "", err
	}

	mirror, err := getFirstServer(out)
	if err != nil {
		return "", err
	}

	return NormalizeMirrorUrl(mirror)
}

func getFirstServer(out []byte) (string, error) {
	servers := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(servers) == 0 {
		return "", errors.New("no server found")
	}

	return servers[0], nil
}

func NormalizeMirrorUrl(mirror string) (string, error) {
	mirrorUrl, err := url.Parse(mirror)
	if err != nil {
		return "", err
	}

	if !mirrorUrl.IsAbs() {
		return "", fmt.Errorf("URL '%s' is not absolute", mirrorUrl.Redacted())
	}

	// strip the core/os/x86_64 path
	return mirrorUrl.JoinPath("../../../").Redacted(), nil
}
