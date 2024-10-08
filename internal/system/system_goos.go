//go:build !linux

package system

import (
	"runtime"
)

func (s *System) GetArchitecture() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return X86_64, nil
	case "arm64":
		return AARCH64, nil
	default:
		return runtime.GOARCH, nil
	}
}
