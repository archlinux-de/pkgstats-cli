package system_test

import (
	"regexp"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetCpuArchitecture(t *testing.T) {
	system := system.System{}

	cpuArch, err := system.GetCpuArchitecture()
	if err != nil {
		t.Error(err)
	}
	if !regexp.MustCompile("^x86_64(_v(2|3|4))?$").MatchString(cpuArch) {
		t.Error(cpuArch)
	}
}
