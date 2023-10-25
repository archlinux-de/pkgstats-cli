package pacman

import (
	"context"
	"errors"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

type Pacman struct {
	timeout    time.Duration
	pacman     string
	pacmanConf string
	repository string
	env        []string
}

func NewPacman() Pacman {
	pacman := Pacman{}
	pacman.timeout = 10 * time.Second
	pacman.pacman = "pacman"
	pacman.pacmanConf = "pacman-conf"
	pacman.repository = "core"
	return pacman
}

func (pacman *Pacman) GetInstalledPackages() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pacman.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, pacman.pacman, "-Qq")
	cmd.Env = pacman.env
	out, err := cmd.Output()

	return strings.Split(strings.TrimSpace(string(out)), "\n"), err
}

func (pacman *Pacman) GetServer() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pacman.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, pacman.pacmanConf, "--repo", pacman.repository, "Server")
	cmd.Env = pacman.env

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
	directories := strings.Split(input, "/")
	path := ""
	if len(directories) > 3 {
		path = strings.Join(directories[:len(directories)-3], "/")
	}
	path = path + "/"

	return path
}
