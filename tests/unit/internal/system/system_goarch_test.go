//go:build !386 && !amd64 && !arm && !arm64

package system_test

import (
	"runtime"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetCpuArchitecture(t *testing.T) {
	system := system.System{}

	cpuArch, err := system.GetCpuArchitecture()
	if err != nil {
		t.Error(err)
	}
	if cpuArch != runtime.GOARCH {
		t.Error(cpuArch)
	}
}
