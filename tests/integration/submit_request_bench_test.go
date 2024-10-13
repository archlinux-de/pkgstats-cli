package integration_test

import (
	"testing"

	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/pacman"
	"pkgstats-cli/internal/system"
)

func BenchmarkCreateRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p, err := pacman.Parse("/etc/pacman.conf")
		if err != nil {
			b.Errorf("Parsing pacman.conf failed: %v", err)
		}
		request, err := submit.CreateRequest(p, system.NewSystem())
		if err != nil {
			b.Errorf("CreateRequest failed: %v", err)
		}
		_ = request // Use the request to avoid compiler optimizations
	}
}
