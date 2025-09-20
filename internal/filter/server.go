package filter

import (
	"net/url"
	"path"
	"strings"
)

func IsFilteredMirrorUrl(filter []string, mirror string) (bool, error) {
	if len(mirror) == 0 {
		return false, nil
	}

	mirrorURL, err := url.Parse(mirror)
	if err != nil {
		return false, err
	}
	mirrorURLHostname := mirrorURL.Hostname()

	for _, filter := range filter {
		isMatch, err := matchMirror(filter, mirrorURLHostname)
		if err != nil {
			return false, err
		}
		if isMatch {
			return true, nil
		}
	}
	return false, nil
}

func matchMirror(pattern string, mirrorURLHostname string) (bool, error) {
	if len(pattern) == 0 {
		return false, nil
	}

	if strings.ContainsAny(pattern, "*?[]") {
		return path.Match(pattern, mirrorURLHostname)
	}

	if patternURL, err := url.Parse(pattern); err == nil && patternURL.Hostname() != "" {
		return patternURL.Hostname() == mirrorURLHostname, nil
	}

	return pattern == mirrorURLHostname, nil
}
