package pkgstats

import (
	"net/http"
	"net/http/httptest"
	"pkgstats-cli/internal/build"
	"testing"
)

func TestSendRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/post" {
			t.Error("/post was not called")
		}
		if req.Method != "POST" {
			t.Error("Invalid Method", req.Method)
		}
		if req.Header.Get("Accept") != "text/plain" {
			t.Error("Invalid Accept Header", req.Header.Get("Accept"))
		}
		if req.UserAgent() != "pkgstats/"+build.Version {
			t.Error("Invalid User-Agent", req.UserAgent())
		}

		if req.PostFormValue("packages") != "pacman\nlinux" {
			t.Error("Invalid packages value", req.FormValue("packages"))
		}
		if req.PostFormValue("cpuarch") != "i686" {
			t.Error("Invalid cpuarch value", req.FormValue("cpuarch"))
		}
		if req.PostFormValue("arch") != "x86_64" {
			t.Error("Invalid arch value", req.FormValue("arch"))
		}
		if req.PostFormValue("mirror") != "https://mirror.pkgbuild.com/" {
			t.Error("Invalid mirror value", req.FormValue("mirror"))
		}
		if req.PostFormValue("quiet") != "false" {
			t.Error("Invalid quiet value", req.FormValue("quiet"))
		}

		rw.Write([]byte("OK"))
	}))
	defer server.Close()

	client := Client{server.Client(), server.URL}
	response, err := client.SendRequest(
		"pacman\nlinux",
		"i686",
		"x86_64",
		"https://mirror.pkgbuild.com/",
		false,
	)

	if err != nil {
		t.Error(err)
	}

	if response != "OK" {
		t.Error("Got wrong response", response)
	}
}
