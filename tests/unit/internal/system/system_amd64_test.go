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
	if !slices.Contains([]string{"x86_64", "x86_64_v2", "x86_64_v3", "x86_64_v4"}, cpuArch) {
		t.Error(cpuArch)
	}
}
