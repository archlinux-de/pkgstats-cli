//go:build !linux

package system_test

import (
	"runtime"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetArchitecture(t *testing.T) {
	cpuArch, err := system.NewSystem().GetCpuArchitecture()
	if err != nil {
		t.Error(err)
	}
	if cpuArch != runtime.GOARCH {
		t.Error(cpuArch)
	}
}

func TestGetOSID(t *testing.T) {
	osId, err := system.NewSystem().GetOSID("")
	if err != nil {
		t.Error(err)
	}
	if osId != runtime.GOOS {
		t.Error(osId)
	}
}
