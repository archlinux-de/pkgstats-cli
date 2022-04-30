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

type result[T any] struct {
	res T
	err error
}

func NewRequest() *Request {
	return &Request{
		Version: "3",
	}
}

func CreateRequest() (*Request, error) {
	pacman := pacman.NewPacman()
	packageChannel := async(pacman.GetInstalledPackages)
	mirrorChannel := async(pacman.GetServer)

	system := system.NewSystem()
	cpuArchitectureChannel := async(system.GetCpuArchitecture)
	architectureChannel := async(system.GetArchitecture)

	packages := <-packageChannel
	if packages.err != nil {
		return nil, packages.err
	}
	mirror := <-mirrorChannel
	if mirror.err != nil {
		return nil, mirror.err
	}
	cpuArchitecture := <-cpuArchitectureChannel
	if cpuArchitecture.err != nil {
		return nil, cpuArchitecture.err
	}
	architecture := <-architectureChannel
	if architecture.err != nil {
		return nil, architecture.err
	}

	request := NewRequest()
	request.System.Architecture = cpuArchitecture.res
	request.OS.Architecture = architecture.res
	request.Pacman.Mirror = mirror.res
	request.Pacman.Packages = packages.res

	return request, nil
}

func async[T any](f func() (T, error)) chan result[T] {
	c := make(chan result[T])
	go func() {
		v, e := f()
		c <- result[T]{v, e}
	}()
	return c
}
