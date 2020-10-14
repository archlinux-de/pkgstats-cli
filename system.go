package main

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

type System struct {
}

func NewSystem() System {
	system := System{}
	return system
}

func (system *System) GetArchitecture() (string, error) {
	out, err := exec.Command("/usr/bin/uname", "-m").Output()
	return strings.TrimSpace(string(out)), err
}

func (system *System) GetCpuArchitecture() (string, error) {
	dat, err := ioutil.ReadFile("/proc/cpuinfo")

	if err == nil && regexp.MustCompile(`(?m)^flags\s*:.*\slm\s`).Match(dat) {
		return "x86_64", nil
	}
	return "", err
}
