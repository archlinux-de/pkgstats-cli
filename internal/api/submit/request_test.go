package submit_test

import (
	"errors"
	"testing"

	"pkgstats-cli/internal/api/submit"
)

type mockPackageManager struct {
	packages []string
	server   string
	pkgErr   error
	srvErr   error
}

func (m *mockPackageManager) GetInstalledPackages() ([]string, error) {
	return m.packages, m.pkgErr
}

func (m *mockPackageManager) GetServer() (string, error) {
	return m.server, m.srvErr
}

type mockSystemInfo struct {
	cpuArch string
	arch    string
	osId    string
	cpuErr  error
	archErr error
	osErr   error
}

func (m *mockSystemInfo) GetCpuArchitecture() (string, error) {
	return m.cpuArch, m.cpuErr
}

func (m *mockSystemInfo) GetArchitecture() (string, error) {
	return m.arch, m.archErr
}

func (m *mockSystemInfo) GetOSId() (string, error) {
	return m.osId, m.osErr
}

func TestCreateRequest(t *testing.T) {
	pm := &mockPackageManager{
		packages: []string{"pacman", "linux"},
		server:   "https://mirror.example.com/",
	}
	si := &mockSystemInfo{
		cpuArch: "x86_64",
		arch:    "x86_64",
		osId:    "linux",
	}

	req, err := submit.CreateRequest(pm, si)
	if err != nil {
		t.Fatal(err)
	}

	if req.Version != submit.Version {
		t.Errorf("Expected version %s, got %s", submit.Version, req.Version)
	}
	if req.System.Architecture != "x86_64" {
		t.Errorf("Expected system architecture x86_64, got %s", req.System.Architecture)
	}
	if req.OS.Architecture != "x86_64" {
		t.Errorf("Expected OS architecture x86_64, got %s", req.OS.Architecture)
	}
	if req.OS.Id != "linux" {
		t.Errorf("Expected OS id linux, got %s", req.OS.Id)
	}
	if req.Pacman.Mirror != "https://mirror.example.com/" {
		t.Errorf("Expected mirror https://mirror.example.com/, got %s", req.Pacman.Mirror)
	}
	if len(req.Pacman.Packages) != 2 {
		t.Errorf("Expected 2 packages, got %d", len(req.Pacman.Packages))
	}
}

func TestCreateRequestReturnsErrorOnPackageFailure(t *testing.T) {
	pm := &mockPackageManager{pkgErr: errors.New("pacman error")}
	si := &mockSystemInfo{}

	_, err := submit.CreateRequest(pm, si)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCreateRequestReturnsErrorOnServerFailure(t *testing.T) {
	pm := &mockPackageManager{packages: []string{"pacman"}, srvErr: errors.New("server error")}
	si := &mockSystemInfo{}

	_, err := submit.CreateRequest(pm, si)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCreateRequestReturnsErrorOnCpuArchFailure(t *testing.T) {
	pm := &mockPackageManager{packages: []string{"pacman"}, server: "https://mirror.example.com/"}
	si := &mockSystemInfo{cpuErr: errors.New("cpu error")}

	_, err := submit.CreateRequest(pm, si)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCreateRequestReturnsErrorOnArchFailure(t *testing.T) {
	pm := &mockPackageManager{packages: []string{"pacman"}, server: "https://mirror.example.com/"}
	si := &mockSystemInfo{cpuArch: "x86_64", archErr: errors.New("arch error")}

	_, err := submit.CreateRequest(pm, si)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCreateRequestReturnsErrorOnOSIdFailure(t *testing.T) {
	pm := &mockPackageManager{packages: []string{"pacman"}, server: "https://mirror.example.com/"}
	si := &mockSystemInfo{cpuArch: "x86_64", arch: "x86_64", osErr: errors.New("os error")}

	_, err := submit.CreateRequest(pm, si)
	if err == nil {
		t.Fatal("Expected error")
	}
}
