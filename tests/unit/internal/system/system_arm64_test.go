package system_test

import (
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetCpuArchitecture(t *testing.T) {
	cpuArch, err := system.NewSystem().GetCpuArchitecture()
	if err != nil {
		t.Error(err)
	}
	if cpuArch != "aarch64" {
		t.Error(cpuArch)
	}
}
