//go:build !linux

package system

import (
	"runtime"
	"testing"
)

func TestGetArchitecture(t *testing.T) {
	system := System{}

	cpuArch, err := system.GetArchitecture()
	if err != nil {
		t.Error(err)
	}
	if cpuArch != runtime.GOARCH {
		t.Error(cpuArch)
	}
}
