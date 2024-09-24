//go:build !386 && !amd64 && !arm && !arm64

package system

import (
	"runtime"
	"testing"
)

func TestGetCpuArchitecture(t *testing.T) {
	system := System{}

	cpuArch, err := system.GetCpuArchitecture()
	if err != nil {
		t.Error(err)
	}
	if cpuArch != runtime.GOARCH {
		t.Error(cpuArch)
	}
}
