package filter

import (
	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/config"
)

func FilterRequest(conf *config.Config, req *submit.Request) error {
	err := filterRequestMirror(conf, req)
	if err != nil {
		return err
	}
	err = filterPackages(conf, req)

	return err
}

func filterRequestMirror(conf *config.Config, req *submit.Request) error {
	isFilteredMirror, err := IsFilteredMirrorUrl(conf.Blocklist.Mirrors, req.Pacman.Mirror)
	if err != nil {
		return err
	}
	if isFilteredMirror {
		req.Pacman.Mirror = ""
	}
	return nil
}

func filterPackages(conf *config.Config, req *submit.Request) error {
	j := 0
	for _, pkg := range req.Pacman.Packages {
		isFilteredPackage, err := IsFilteredPackage(conf.Blocklist.Packages, pkg)
		if err != nil {
			return err
		}

		if !isFilteredPackage {
			req.Pacman.Packages[j] = pkg
			j++
		}
	}

	req.Pacman.Packages = req.Pacman.Packages[:j]

	return nil
}
