package filter_test

import (
	"testing"

	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/config"
	"pkgstats-cli/internal/filter"
)

func TestFilterRequestRemovesBlockedPackages(t *testing.T) {
	conf := &config.Config{}
	conf.Blocklist.Packages = []string{"secret-*"}

	req := &submit.Request{
		Pacman: submit.Pacman{
			Packages: []string{"pacman", "secret-pkg", "linux"},
			Mirror:   "http://example.com/",
		},
	}

	if err := filter.FilterRequest(conf, req); err != nil {
		t.Fatal(err)
	}

	if len(req.Pacman.Packages) != 2 {
		t.Fatalf("Expected 2 packages, got %d: %v", len(req.Pacman.Packages), req.Pacman.Packages)
	}
	for _, pkg := range req.Pacman.Packages {
		if pkg == "secret-pkg" {
			t.Error("Expected secret-pkg to be filtered")
		}
	}
}

func TestFilterRequestClearsBlockedMirror(t *testing.T) {
	conf := &config.Config{}
	conf.Blocklist.Mirrors = []string{"private.lan"}

	req := &submit.Request{
		Pacman: submit.Pacman{
			Packages: []string{"pacman"},
			Mirror:   "http://private.lan/archlinux",
		},
	}

	if err := filter.FilterRequest(conf, req); err != nil {
		t.Fatal(err)
	}

	if req.Pacman.Mirror != "" {
		t.Errorf("Expected mirror to be cleared, got %s", req.Pacman.Mirror)
	}
}

func TestFilterRequestKeepsUnblockedMirror(t *testing.T) {
	conf := &config.Config{}
	conf.Blocklist.Mirrors = []string{"private.lan"}

	mirror := "http://public.example.com/archlinux"
	req := &submit.Request{
		Pacman: submit.Pacman{
			Packages: []string{"pacman"},
			Mirror:   mirror,
		},
	}

	if err := filter.FilterRequest(conf, req); err != nil {
		t.Fatal(err)
	}

	if req.Pacman.Mirror != mirror {
		t.Errorf("Expected mirror to remain %s, got %s", mirror, req.Pacman.Mirror)
	}
}

func TestFilterRequestWithEmptyConfig(t *testing.T) {
	conf := &config.Config{}
	req := &submit.Request{
		Pacman: submit.Pacman{
			Packages: []string{"pacman", "linux"},
			Mirror:   "http://example.com/",
		},
	}

	if err := filter.FilterRequest(conf, req); err != nil {
		t.Fatal(err)
	}

	if len(req.Pacman.Packages) != 2 {
		t.Errorf("Expected 2 packages, got %d", len(req.Pacman.Packages))
	}
	if req.Pacman.Mirror != "http://example.com/" {
		t.Errorf("Expected mirror unchanged, got %s", req.Pacman.Mirror)
	}
}

func TestFilterRequestReturnsErrorOnInvalidPackagePattern(t *testing.T) {
	conf := &config.Config{}
	conf.Blocklist.Packages = []string{"["}

	req := &submit.Request{
		Pacman: submit.Pacman{
			Packages: []string{"pacman"},
			Mirror:   "http://example.com/",
		},
	}

	if err := filter.FilterRequest(conf, req); err == nil {
		t.Error("Expected error for invalid glob pattern")
	}
}

func TestFilterRequestReturnsErrorOnInvalidMirrorPattern(t *testing.T) {
	conf := &config.Config{}
	conf.Blocklist.Mirrors = []string{"["}

	req := &submit.Request{
		Pacman: submit.Pacman{
			Packages: []string{"pacman"},
			Mirror:   "http://example.com/",
		},
	}

	if err := filter.FilterRequest(conf, req); err == nil {
		t.Error("Expected error for invalid glob pattern")
	}
}
