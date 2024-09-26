package pacman

import (
	"context"
	"errors"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const timeout = 10 * time.Second

type Pacman struct {
	Timeout    time.Duration
	Pacman     string
	PacmanConf string
	Repository string
	Env        []string
}

func NewPacman() Pacman {
	return Pacman{
		Timeout:    timeout,
		Pacman:     "pacman",
		PacmanConf: "pacman-conf",
		Repository: "core",
	}
}

func (pacman *Pacman) GetInstalledPackages() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pacman.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, pacman.Pacman, "-Qq")
	cmd.Env = pacman.Env
	out, err := cmd.Output()

	return strings.Split(strings.TrimSpace(string(out)), "\n"), err
}

func (pacman *Pacman) GetServer() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pacman.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, pacman.PacmanConf, "--repo", pacman.Repository, "Server")
	cmd.Env = pacman.Env

	out, err := cmd.Output()

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
