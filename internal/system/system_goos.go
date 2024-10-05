//go:build !linux

package system

import (
	"runtime"
)

func (system *System) GetArchitecture() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "x86_64", nil
	case "arm64":
		return "aarch64", nil
	default:
		return runtime.GOARCH, nil
	}
}
