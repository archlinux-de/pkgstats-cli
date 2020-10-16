package pacman

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
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
	pacman.pacman = "/usr/bin/pacman"
	pacman.pacmanConf = "/usr/bin/pacman-conf"
	pacman.repository = "core"
	return pacman
}

func (pacman *Pacman) GetInstalledPackages() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pacman.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, pacman.pacman, "-Qq")
	cmd.Env = pacman.env
	out, err := cmd.Output()

	return strings.TrimSpace(string(out)), err
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
		return mirror, errors.New("No server found")
	}

	mirrorUrl, _ := url.Parse(mirror)
	path := regexp.MustCompile(fmt.Sprintf(`^(.*/)%s/os/.*`, pacman.repository)).ReplaceAllString(mirrorUrl.Path, "$1")

	port := ""
	if mirrorUrl.Port() != "" {
		port = ":" + mirrorUrl.Port()
	}

	return mirrorUrl.Scheme + "://" + mirrorUrl.Hostname() + port + path, err
}
