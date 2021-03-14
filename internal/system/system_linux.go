package system

import (
	"bytes"
	"golang.org/x/sys/unix"
)

func (system *System) GetArchitecture() (string, error) {
	var utsname unix.Utsname
	err := unix.Uname(&utsname)
	if err != nil {
		return "", err
	}

	return string(utsname.Machine[:bytes.IndexByte(utsname.Machine[:], 0)]), nil
}
