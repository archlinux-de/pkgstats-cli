package system

import (
	"bytes"
	"os"
	"runtime"
	"strings"

	"golang.org/x/sys/unix"
)

func (s *System) GetArchitecture() (string, error) {
	var utsname unix.Utsname
	err := unix.Uname(&utsname)
	if err != nil {
		return "", err
	}

	return string(utsname.Machine[:bytes.IndexByte(utsname.Machine[:], 0)]), nil
}

func (s *System) GetOSID(osReleasePaths ...string) (string, error) {
	var id string

	if len(osReleasePaths) == 0 {
		osReleasePaths = []string{"/etc/os-release", "/usr/lib/os-release"}
	}

	for _, osReleasePath := range osReleasePaths {
		content, err := os.ReadFile(osReleasePath)
		if err != nil {
			continue
		}

		lines := strings.Split(string(content), "\n")
		const keyValueParts = 2
		for _, line := range lines {
			if trimmedLine := strings.TrimSpace(line); trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", keyValueParts)
			if len(parts) != keyValueParts {
				continue
			}

			key := strings.TrimSpace(parts[0])
			if key == "ID" {
				id = strings.Trim(strings.TrimSpace(parts[1]), `"'`)
			}
		}
		break
	}

	if id == "" {
		id = runtime.GOOS
	}

	return id, nil
}
