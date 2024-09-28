package integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"slices"

	"pkgstats-cli/internal/api/request"
	"pkgstats-cli/internal/api/submit"
	"pkgstats-cli/internal/system"
)

func NewServer() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/submit", handleSubmit)
	mux.HandleFunc("GET /api/packages", handlePackages)
	mux.HandleFunc("GET /api/packages/pacman", handlePackagesPacman)
	mux.HandleFunc("GET /api/packages/php", handlePackagesPhp)
	mux.HandleFunc("/", handleDefault)
	return mux
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if !regexp.MustCompile(`^pkgstats/[\w.-]+$`).MatchString(r.Header.Get("User-Agent")) {
		http.Error(w, "Invalid user agent", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("redirect") == "" {
		w.Header().Set("Location", "/api/submit?redirect=1")
		w.WriteHeader(http.StatusPermanentRedirect)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request submit.Request
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if isValidSubmitRequest(&request) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "TEST FAILED", http.StatusBadRequest)
	}
}

func handlePackages(w http.ResponseWriter, r *http.Request) {
	response := request.PackagePopularityList{
		Total: 42,
		Count: 2,
		PackagePopularities: []request.PackagePopularity{
			{Name: "php", Popularity: 56.78},
			{Name: "php-fpm", Popularity: 12.34},
		},
	}
	_ = json.NewEncoder(w).Encode(response)
}

func handlePackagesPacman(w http.ResponseWriter, r *http.Request) {
	response := request.PackagePopularity{Name: "pacman", Popularity: 12.34}
	_ = json.NewEncoder(w).Encode(response)
}

func handlePackagesPhp(w http.ResponseWriter, r *http.Request) {
	response := request.PackagePopularity{Name: "php", Popularity: 56.78}
	_ = json.NewEncoder(w).Encode(response)
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unknown request", http.StatusBadRequest)
}

func isValidSubmitRequest(request *submit.Request) bool {
	system := system.NewSystem()
	osArchitecture, _ := system.GetArchitecture()
	cpuArchitecture, _ := system.GetCpuArchitecture()

	return request.Version == "3" &&
		request.OS.Architecture == osArchitecture &&
		request.System.Architecture == cpuArchitecture &&
		regexp.MustCompile(`^https?://.+$`).MatchString(request.Pacman.Mirror) &&
		len(request.Pacman.Packages) > 1 &&
		slices.Contains(request.Pacman.Packages, "pacman-mirrorlist")
}
