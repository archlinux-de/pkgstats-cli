package integration_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pkgstats-cli/internal/system"
	"regexp"
	"slices"
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
	userAgent := r.Header.Get("User-Agent")
	requestUri := r.RequestURI
	log.Printf("Got request from %s on %s", userAgent, requestUri)

	if !regexp.MustCompile(`^pkgstats/[\w.-]+$`).MatchString(userAgent) {
		http.Error(w, "Invalid user agent", http.StatusBadRequest)
		log.Println("Invalid user agent")
		return
	}

	if r.URL.Query().Get("redirect") == "" {
		log.Println("Testing redirect")
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
		log.Println("Request was valid")
		w.WriteHeader(http.StatusNoContent)
	} else {
		log.Println("Request was invalid")
		log.Println(request)
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
	json.NewEncoder(w).Encode(response)
}

func handlePackagesPacman(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"name":       "pacman",
		"popularity": 12.34,
	}
	json.NewEncoder(w).Encode(response)
}

func handlePackagesPhp(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"name":       "php",
		"popularity": 56.78,
	}
	json.NewEncoder(w).Encode(response)
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	log.Println("Unknown request")
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
