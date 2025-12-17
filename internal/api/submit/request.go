package submit

type System struct {
	Architecture string `json:"architecture"`
}

type OS struct {
	Architecture string `json:"architecture"`
	Id           string `json:"id"`
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

type PackageManager interface {
	GetInstalledPackages() ([]string, error)
	GetServer() (string, error)
}

type SystemInfo interface {
	GetCpuArchitecture() (string, error)
	GetArchitecture() (string, error)
	GetOSId() (string, error)
}

func CreateRequest(p PackageManager, s SystemInfo) (*Request, error) {
	packages, err := p.GetInstalledPackages()
	if err != nil {
		return nil, err
	}
	mirror, err := p.GetServer()
	if err != nil {
		return nil, err
	}

	cpuArchitecture, err := s.GetCpuArchitecture()
	if err != nil {
		return nil, err
	}
	architecture, err := s.GetArchitecture()
	if err != nil {
		return nil, err
	}

	osId, err := s.GetOSId()
	if err != nil {
		return nil, err
	}

	return &Request{
		Version: Version,
		System:  System{Architecture: cpuArchitecture},
		OS:      OS{Architecture: architecture, Id: osId},
		Pacman:  Pacman{Packages: packages, Mirror: mirror},
	}, nil
}
