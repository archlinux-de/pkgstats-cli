package submit

import (
	"pkgstats-cli/internal/pacman"
	"pkgstats-cli/internal/system"
)

type System struct {
	Architecture string `json:"architecture"`
}

type OS struct {
	Architecture string `json:"architecture"`
}

type Pacman struct {
	Mirror   string   `json:"mirror"`
	Packages []string `json:"packages"`
}

type Request struct {
	Version string `json:"version"`
	System  System `json:"system"`
	OS      OS     `json:"os"`
	Pacman  Pacman `json:"pacman"`
}

const Version = "3"

func CreateRequest() (*Request, error) {
	p := pacman.NewPacman()
	packages, err := p.GetInstalledPackages()
	if err != nil {
		return nil, err
	}
	mirror, err := p.GetServer()
	if err != nil {
		return nil, err
	}

	s := system.NewSystem()
	cpuArchitecture, err := s.GetCpuArchitecture()
	if err != nil {
		return nil, err
	}
	architecture, err := s.GetArchitecture()
	if err != nil {
		return nil, err
	}

	return &Request{
		Version: Version,
		System:  System{Architecture: cpuArchitecture},
		OS:      OS{Architecture: architecture},
		Pacman:  Pacman{Packages: packages, Mirror: mirror},
	}, nil
}
