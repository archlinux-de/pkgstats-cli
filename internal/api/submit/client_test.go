package submit

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"pkgstats-cli/internal/build"
	"strings"
	"testing"
)

func TestSendRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/api/submit" {
			t.Error("/api/submit was not called")
		}

		validateRequest(t, req)

		rw.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := Client{server.Client(), server.URL}
	request := NewRequest()
	request.Pacman.Packages = []string{"pacman", "linux"}
	request.System.Architecture = "x86_64"
	request.OS.Architecture = "i686"
	request.Pacman.Mirror = "https://geo.mirror.pkgbuild.com/"
	err := client.SendRequest(*request)

	if err != nil {
		t.Error(err)
	}
}

func validateRequest(t *testing.T, req *http.Request) {
	if req.Method != http.MethodPost {
		t.Error("Invalid Method", req.Method)
	}
	if req.Header.Get("Accept") != "application/json" {
		t.Error("Invalid Accept Header", req.Header.Get("Accept"))
	}
	if req.UserAgent() != "pkgstats/"+build.Version {
		t.Error("Invalid User-Agent", req.UserAgent())
	}

	request := Request{}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		t.Errorf("Could not read request body %s", err)
	}
	err = json.Unmarshal(body, &request)
	if err != nil {
		t.Errorf("Could not decode JSON: %s", err)
	}

	if request.Version != "3" {
		t.Error("Invalid version value")
	}
	if strings.Join(request.Pacman.Packages, ",") != "pacman,linux" {
		t.Error("Invalid packages value")
	}
	if request.System.Architecture != "x86_64" {
		t.Error("Invalid cpuarch value")
	}
	if request.OS.Architecture != "i686" {
		t.Error("Invalid arch value")
	}
	if request.Pacman.Mirror != "https://geo.mirror.pkgbuild.com/" {
		t.Error("Invalid mirror value")
	}
}

func TestSendRequestFollowsRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("redirect") == "" {
			rw.Header().Add("Location", "/api/submit?redirect=1")
			rw.WriteHeader(http.StatusPermanentRedirect)

			return
		}
		if req.URL.String() != "/api/submit?redirect=1" {
			t.Error("/api/submit?redirect=1 was not called")
		}

		validateRequest(t, req)

		rw.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := Client{server.Client(), server.URL}
	request := NewRequest()
	request.Pacman.Packages = []string{"pacman", "linux"}
	request.System.Architecture = "x86_64"
	request.OS.Architecture = "i686"
	request.Pacman.Mirror = "https://geo.mirror.pkgbuild.com/"
	err := client.SendRequest(*request)

	if err != nil {
		t.Error(err)
	}
}
