package system

import (
	"regexp"
	"testing"
)

func TestGetCpuArchitecture(t *testing.T) {
	system := System{}

	cpuArch, err := system.GetCpuArchitecture()

	if err != nil {
		t.Error(err)
	}
	if cpuArch == "aarch64" {
		t.Error(cpuArch)
	}
}
