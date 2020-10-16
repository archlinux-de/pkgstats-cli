package build

import "testing"

func TestVersion(t *testing.T) {
	if Version != "dev" {
		t.Error("Version is not defined")
	}
}
