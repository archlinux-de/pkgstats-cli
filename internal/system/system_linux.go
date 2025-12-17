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

func (s *System) GetOSId() (string, error) {
	for _, path := range []string{"/etc/os-release", "/usr/lib/os-release"} {
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if id := ParseOSId(content); id != "" {
			return id, nil
		}
		break
	}

	return runtime.GOOS, nil
}

func ParseOSId(content []byte) string {
	var id string

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

	return id
}
