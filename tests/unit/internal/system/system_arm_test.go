package system_test

import (
	"slices"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetCpuArchitecture(t *testing.T) {
	cpuArch, err := system.NewSystem().GetCpuArchitecture()
	if err != nil {
		t.Error(err)
	}
	if !slices.Contains([]string{system.ARMV5, system.ARMV6, system.ARMV7, system.AARCH64}, cpuArch) {
		t.Error(cpuArch)
	}
}
