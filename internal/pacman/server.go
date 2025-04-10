package pacman

import (
	"fmt"
	"net/url"
)

func (p *Pacman) GetServer() (string, error) {
	return NormalizeMirrorUrl(p.config.server)
}

func NormalizeMirrorUrl(mirror string) (string, error) {
	mirrorUrl, err := url.Parse(mirror)
	if err != nil {
		return "", err
	}

	if !mirrorUrl.IsAbs() {
		return "", fmt.Errorf("URL '%s' is not absolute", mirrorUrl.Redacted())
	}

	// remove login information
	mirrorUrl.User = nil

	// strip the core/os/x86_64 path
	return mirrorUrl.JoinPath("../../../").Redacted(), nil
}
