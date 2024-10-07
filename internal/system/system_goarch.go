//go:build !386 && !amd64 && !arm && !arm64

package system

import (
	"runtime"
)

func (s *System) GetCpuArchitecture() (string, error) {
	return runtime.GOARCH, nil
}
