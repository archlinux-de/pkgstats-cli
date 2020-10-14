package main

import (
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

type Pacman struct {
}

func NewPacman() Pacman {
	pacman := Pacman{}
	return pacman
}

func (pacman *Pacman) GetInstalledPackages() (string, error) {
	out, err := exec.Command("/usr/bin/pacman", "-Qq").Output()
	return strings.TrimSpace(string(out)), err
}

func (pacman *Pacman) GetServer() (string, error) {
	out, err := exec.Command("/usr/bin/pacman-conf", "--repo", "extra", "Server").Output()
	mirror := strings.TrimSpace(string(out))
	url, _ := url.Parse(mirror)
	path := regexp.MustCompile(`^(.*/)extra/os/.*`).ReplaceAllString(url.Path, "$1")

	port := ""
	if url.Port() != "" {
		port = ":" + url.Port()
	}

	return url.Scheme + "://" + url.Hostname() + port + path, err
}
