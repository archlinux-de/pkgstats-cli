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
	if !regexp.MustCompile("^(i686|x86_64(_v(2|3|4))?)$").MatchString(cpuArch) {
		t.Error(cpuArch)
	}
}
