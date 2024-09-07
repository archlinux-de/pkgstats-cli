package system

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
)

func TestGetMachine(t *testing.T) {
	system := System{}

	arch, err := system.GetArchitecture()

	expectedArch := fmt.Sprintf("^%s$", runtime.GOARCH)
	switch runtime.GOARCH {
	case "amd64":
		expectedArch = "^x86_64$"
	case "386":
		expectedArch = "^i686$"
	case "arm":
		expectedArch = "^armv(5|6|7)"
	case "arm64":
		expectedArch = "^aarch64$"
	case "loong64":
		expectedArch = "^loongarch64$"
	}

	if err != nil {
		t.Error(err)
	}
	if !regexp.MustCompile(expectedArch).MatchString(arch) {
		t.Error(arch)
	}
}
