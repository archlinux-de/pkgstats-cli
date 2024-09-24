package integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"slices"

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

	var request map[string]interface{}
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if isValidRequest(request) {
		w.WriteHeader(http.StatusNoContent)
	} else {
		http.Error(w, "TEST FAILED", http.StatusBadRequest)
	}
}

func handlePackages(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"total": 42,
		"count": 2,
		"packagePopularities": []map[string]interface{}{
			{"name": "php", "popularity": 56.78},
			{"name": "php-fpm", "popularity": 12.34},
		},
	}
	_ = json.NewEncoder(w).Encode(response)
}

func handlePackagesPacman(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"name":       "pacman",
		"popularity": 12.34,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func handlePackagesPhp(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"name":       "php",
		"popularity": 56.78,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Unknown request", http.StatusBadRequest)
}

func isValidRequest(request map[string]interface{}) bool {
	system := system.NewSystem()
	osArchitecture, _ := system.GetArchitecture()
	cpuArchitecture, _ := system.GetCpuArchitecture()

	return request["version"] == "3" &&
		request["os"].(map[string]interface{})["architecture"] == osArchitecture &&
		request["system"].(map[string]interface{})["architecture"] == cpuArchitecture &&
		regexp.MustCompile(`^https?://.+$`).MatchString(request["pacman"].(map[string]interface{})["mirror"].(string)) &&
		len(request["pacman"].(map[string]interface{})["packages"].([]interface{})) > 1 &&
		slices.Contains(request["pacman"].(map[string]interface{})["packages"].([]interface{}), "pacman-mirrorlist")
}
