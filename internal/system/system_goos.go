//go:build !linux

package system

import (
	"runtime"
)

func (system *System) GetArchitecture() (string, error) {
	return runtime.GOARCH, nil
}
