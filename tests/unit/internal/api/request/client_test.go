package request_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pkgstats-cli/internal/api/request"
)

func TestGetPackages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.String(), "/api/packages/") {
			t.Error("/api/packages/ was not called")
		}
		fmt.Fprint(rw, "{}")
	}))
	defer server.Close()

	client := request.Client{Client: server.Client(), BaseURL: server.URL}
	_, errs := client.GetPackages("foo", "bar")
	if len(errs) > 0 {
		t.Error(errs)
	}
}

func TestGetPackage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/api/packages/foo" {
			t.Error("/api/packages/foo was not called")
		}
		fmt.Fprint(rw, "{}")
	}))
	defer server.Close()

	client := request.Client{Client: server.Client(), BaseURL: server.URL}
	_, err := client.GetPackage("foo")
	if err != nil {
		t.Error(err)
	}
}

func TestSearchPackages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/api/packages?limit=1&query=foo" {
			t.Error("/api/packages?limit=1&query=foo was not called")
		}
		fmt.Fprint(rw, "{}")
	}))
	defer server.Close()

	client := request.Client{Client: server.Client(), BaseURL: server.URL}
	_, err := client.SearchPackages("foo", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestHandleServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "Ignored server message")
	}))
	defer server.Close()

	client := request.Client{Client: server.Client(), BaseURL: server.URL}
	_, errs := client.GetPackages("foo")
	if len(errs) == 0 {
		t.Error("Expected error, got none")
	} else if errs[0].Error() != "Bad Request" {
		t.Errorf("Expected Bad Request, got %v", errs[0])
	}
}
