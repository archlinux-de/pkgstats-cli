//go:build !linux

package system_test

import (
	"runtime"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetArchitecture(t *testing.T) {
	system := system.System{}

	cpuArch, err := system.GetArchitecture()
	if err != nil {
		t.Error(err)
	}
	if cpuArch != runtime.GOARCH {
		t.Error(cpuArch)
	}
}
