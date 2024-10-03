package system_test

import (
	"slices"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetCpuArchitecture(t *testing.T) {
	system := system.System{}

	cpuArch, err := system.GetCpuArchitecture()
	if err != nil {
		t.Error(err)
	}
	if !slices.Contains([]string{"armv5", "armv6", "armv7", "aarch64"}, cpuArch) {
		t.Error(cpuArch)
	}
}
