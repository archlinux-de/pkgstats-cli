package pacman

import (
	"context"
	"errors"
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

	servers := strings.Split(strings.TrimSpace(string(out)), "\n")
	mirror := ""
	if len(servers) > 0 {
		mirror = servers[0]
	} else {
		return mirror, errors.New("no server found")
	}

	mirrorUrl, _ := url.Parse(mirror)
	path := extractMirrorPath(mirrorUrl.Path)
	port := extractMirrorPort(mirrorUrl.Port())

	return mirrorUrl.Scheme + "://" + mirrorUrl.Hostname() + port + path, err
}

func extractMirrorPort(input string) string {
	port := ""
	if input != "" {
		port = ":" + input
	}

	return port
}

func extractMirrorPath(input string) string {
	const directoryPattern = 3
	directories := strings.Split(input, "/")
	path := ""
	if len(directories) > directoryPattern {
		path = strings.Join(directories[:len(directories)-directoryPattern], "/")
	}
	path += "/"

	return path
}
