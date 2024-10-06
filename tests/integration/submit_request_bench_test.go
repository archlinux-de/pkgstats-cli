package integration_test

import (
	"testing"

	"pkgstats-cli/internal/api/submit"
)

func BenchmarkCreateRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		request, err := submit.CreateRequest()
		if err != nil {
			b.Errorf("CreateRequest failed: %v", err)
		}
		_ = request // Use the request to avoid compiler optimizations
	}
}
