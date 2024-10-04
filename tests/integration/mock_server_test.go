package integration_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"

	"pkgstats-cli/internal/api/request"
	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/system"
)

type PackagePopularityListMock struct {
	request.PackagePopularityList
}

func NewPackagePopularityListMock() *PackagePopularityListMock {
	return &PackagePopularityListMock{
		request.PackagePopularityList{
			Total: 4,
			Count: 3,
			PackagePopularities: []request.PackagePopularity{
				{Name: "php", Popularity: 56.78},
				{Name: "php-fpm", Popularity: 12.34},
				{Name: "pacman", Popularity: 10.78},
			},
		},
	}
}

func (p *PackagePopularityListMock) getByName(name string) (request.PackagePopularity, error) {
	for _, pkg := range p.PackagePopularities {
		if pkg.Name == name {
			return pkg, nil
		}
	}
	return request.PackagePopularity{}, fmt.Errorf("package %s not found", name)
}

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/submit", handleSubmit)
	mux.HandleFunc("GET /api/packages", handlePackages)
	mux.HandleFunc("GET /api/packages/{package}", handlePackage)
	return mux
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if !regexp.MustCompile(`^pkgstats/[\w.-]+$`).MatchString(r.Header.Get("User-Agent")) {
		http.Error(w, fmt.Sprintf("Invalid user agent %s", r.Header.Get("User-Agent")), http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("redirect") == "" {
		w.Header().Set("Location", "/api/submit?redirect=1")
		w.WriteHeader(http.StatusPermanentRedirect)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request submit.Request
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validateSubmitRequest(w, &request)
}

func handlePackages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(NewPackagePopularityListMock()); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func handlePackage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p, err := NewPackagePopularityListMock().getByName(r.PathValue("package"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching package: %v", err), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func validateSubmitRequest(w http.ResponseWriter, request *submit.Request) {
	s := system.NewSystem()
	osArchitecture, err := s.GetArchitecture()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cpuArchitecture, err := s.GetCpuArchitecture()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.Version != "3" {
		http.Error(w, fmt.Sprintf("Expected version 3, got %s", request.Version), http.StatusBadRequest)
		return
	}
	if request.OS.Architecture != osArchitecture {
		http.Error(w, fmt.Sprintf("Expected OS architecture %s, got %s", osArchitecture, request.OS.Architecture), http.StatusBadRequest)
		return
	}
	if request.System.Architecture != cpuArchitecture {
		http.Error(w, fmt.Sprintf("Expected CPU architecture %s, got %s", cpuArchitecture, request.System.Architecture), http.StatusBadRequest)
		return
	}
	if !regexp.MustCompile(`^https?://.+$`).MatchString(request.Pacman.Mirror) {
		http.Error(w, fmt.Sprintf("Invalid HTTP mirror URL: %s", request.Pacman.Mirror), http.StatusBadRequest)
		return
	}
	if len(request.Pacman.Packages) <= 1 {
		http.Error(w, "Expected more than 1 package", http.StatusBadRequest)
		return
	}
	if !slices.Contains(request.Pacman.Packages, "pacman-mirrorlist") {
		http.Error(w, "Expected pacakge list to contain pacman-mirrorlist", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
