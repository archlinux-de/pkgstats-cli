package pacman

import (
	"fmt"
	"os"
	"path/filepath"
)

func (p *Pacman) GetInstalledPackages() ([]string, error) {
	const verParts = 2

	db, err := os.Open(filepath.Join(p.config.dbPath, "local"))
	if err != nil {
		return nil, fmt.Errorf("failed to open pacman DB: %w", err)
	}
	defer db.Close()

	packages, err := db.Readdirnames(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read package names: %w", err)
	}

	// Reuse the memeory of packages
	filteredPackages := packages[:0]
	for _, pkg := range packages {
		hyphenCount := 0
		for j := len(pkg) - 1; j >= 0; j-- {
			if pkg[j] == '-' {
				hyphenCount++
				if hyphenCount == verParts {
					filteredPackages = append(filteredPackages, pkg[:j])
					break
				}
			}
		}
	}

	return filteredPackages, nil
}
