package system_test

import (
	"runtime"
	"slices"
	"strings"
	"testing"

	"pkgstats-cli/internal/system"
)

func TestGetMachine(t *testing.T) {
	arch, err := system.NewSystem().GetArchitecture()

	expectedArch := func(arch string) bool { return arch == runtime.GOARCH }
	switch runtime.GOARCH {
	case "amd64":
		expectedArch = func(arch string) bool { return strings.HasPrefix(arch, system.X86_64) }
	case "386":
		expectedArch = func(arch string) bool { return slices.Contains([]string{system.I586, system.I686}, arch) }
	case "arm":
		expectedArch = func(arch string) bool { return strings.HasPrefix(arch, "armv") }
	case "arm64":
		expectedArch = func(arch string) bool { return arch == system.AARCH64 }
	case "loong64":
		expectedArch = func(arch string) bool { return arch == "loongarch64" }
	}

	if err != nil {
		t.Error(err)
	}
	if !expectedArch(arch) {
		t.Error(arch)
	}
}
