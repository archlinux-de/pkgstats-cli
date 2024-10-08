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
	if !slices.Contains([]string{system.X86_64, system.X86_64_V2, system.X86_64_V3, system.X86_64_V4}, cpuArch) {
		t.Error(cpuArch)
	}
}
