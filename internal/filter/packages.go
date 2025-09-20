package filter

import "path"

func IsFilteredPackage(filters []string, pkg string) (bool, error) {
	for _, filter := range filters {
		isFiltered, err := path.Match(filter, pkg)
		if err != nil {
			return false, err
		}
		if isFiltered {
			return true, nil
		}
	}

	return false, nil
}
