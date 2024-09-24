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
	if !regexp.MustCompile("^(armv(5|6|7)|aarch64)$").MatchString(cpuArch) {
		t.Error(cpuArch)
	}
}
