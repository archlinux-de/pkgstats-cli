package system_test

import (
	"runtime"
	"slices"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetMachine(t *testing.T) {
	system := system.System{}

	arch, err := system.GetArchitecture()

	expectedArch := []string{runtime.GOARCH}
	switch runtime.GOARCH {
	case "amd64":
		expectedArch = []string{"x86_64"}
	case "386":
		expectedArch = []string{"i386", "i486", "i586", "i686"}
	case "arm":
		expectedArch = []string{"arm", "armv4l", "armv5l", "armv5tejl", "armv5tel", "armv6l", "armv7l", "armv7hl", "armv8l"}
	case "arm64":
		expectedArch = []string{"aarch64", "aarch64_be"}
	case "loong64":
		expectedArch = []string{"loongarch64"}
	}

	if err != nil {
		t.Error(err)
	}
	if !slices.Contains(expectedArch, arch) {
		t.Error(arch)
	}
}
