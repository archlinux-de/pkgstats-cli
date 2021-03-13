package system

import (
	"os/exec"
	"strings"
)

type System struct {
	env     []string
	uname   string
}

func NewSystem() System {
	system := System{}
	system.uname = "uname"
	return system
}

func (system *System) GetArchitecture() (string, error) {
	arch, err := system.getMachine()
	return arch, err
}

func (system *System) getMachine() (string, error) {
	cmd := exec.Command(system.uname, "-m")
	cmd.Env = system.env
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
