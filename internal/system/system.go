package system

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

type System struct {
	env     []string
	uname   string
	cpuInfo string
}

func NewSystem() System {
	system := System{}
	system.uname = "/usr/bin/uname"
	system.cpuInfo = "/proc/cpuinfo"
	return system
}

func (system *System) GetArchitecture() (string, error) {
	arch, err := system.getMachine()
	return arch, err
}

func (system *System) GetCpuArchitecture() (string, error) {
	architecture, err := system.GetArchitecture()

	// detect a 64 bit CPU when ruinning a 32 bit OS
	if architecture == "i686" && system.hasLongMode() {
		architecture = "x86_64"
	}
	return architecture, err
}

func (system *System) hasLongMode() bool {
	cpuInfo, err := ioutil.ReadFile(system.cpuInfo)

	return err == nil && regexp.MustCompile(`(?m)^flags\s*:[^\n]*\blm\b[^\n]*$`).Match(cpuInfo)
}

func (system *System) getMachine() (string, error) {
	cmd := exec.Command(system.uname, "-m")
	cmd.Env = system.env
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
