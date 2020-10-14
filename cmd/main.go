package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Version pkgstats version
var Version = "0.0.0-dev"

func main() {
	version := flag.Bool("v", false, "show the version of pkgstats")
	debug := flag.Bool("d", false, "enable debug mode")
	dryRun := flag.Bool("s", false, "show what information would be sent\n(but do not send anything)")
	quiet := flag.Bool("q", false, "be quiet except on errors")
	flag.Usage = printUsage
	flag.Parse()

	if *version {
		fmt.Println("pkgstats, version ", Version)
		os.Exit(0)
	}

	if *debug {
		fmt.Println("Debug mode is not available yet!")
		os.Exit(1)
	}

	if !*quiet {
		fmt.Println("Collecting data...")
	}

	if !*quiet {
		fmt.Println("Submitting data...")
	}
	if !*dryRun {
		response, err := sendRequest(
			getPackages(),
			getCpuArchitecture(),
			getArchitecture(),
			getMirror(),
			*quiet,
		)
		if err != nil {
			fmt.Println("Sorry, data could not be sent.")
			fmt.Println(err)
		}
		if !*quiet {
			fmt.Println(response)
		}
	} else {
		fmt.Println("packages=")
		fmt.Println(getPackages())

		fmt.Println("")

		fmt.Println("arch=", getArchitecture())
		fmt.Println("cpuarch=", getCpuArchitecture())
		fmt.Println("pkgstatsver=", Version)
		fmt.Println("mirror=", getMirror())
		fmt.Println("quiet=", *quiet)
	}
}

func printUsage() {
	path, _ := os.Executable()

	fmt.Printf("usage: %s [option]\n", path)
	fmt.Println("options:")

	flag.PrintDefaults()

	fmt.Println("")
	fmt.Println("pkgstats sends a list of all installed packages,")
	fmt.Println("the architecture and the mirror you are using")
	fmt.Println("to the Arch Linux project.")
	fmt.Println("")
	fmt.Println("Statistics are available at https://pkgstats.archlinux.de/")
}

func getArchitecture() string {
	out, _ := exec.Command("uname", "-m").Output()
	return strings.TrimSpace(string(out))
}

func getCpuArchitecture() string {
	dat, _ := ioutil.ReadFile("/proc/cpuinfo")

	if regexp.MustCompile(`(?m)^flags\s*:.*\slm\s`).Match(dat) {
		return "x86_64"
	}
	return ""
}

func getMirror() string {
	out, _ := exec.Command("pacman-conf", "--repo", "extra", "Server").Output()
	mirror := strings.TrimSpace(string(out))
	url, _ := url.Parse(mirror)
	path := regexp.MustCompile(`^(.*/)extra/os/.*`).ReplaceAllString(url.Path, "$1")

	port := ""
	if url.Port() != "" {
		port = ":" + url.Port()
	}

	return url.Scheme + "://" + url.Hostname() + port + path
}

func getPackages() string {
	out, _ := exec.Command("pacman", "-Qq").Output()
	return strings.TrimSpace(string(out))
}

func sendRequest(packages string, cpuArchitecture string, architecture string, mirror string, quiet bool) (string, error) {
	form := url.Values{}
	form.Add("packages", packages)
	form.Add("arch", architecture)
	form.Add("cpuarch", cpuArchitecture)
	form.Add("mirror", mirror)
	if quiet {
		form.Add("quiet", "true")
	} else {
		form.Add("quiet", "false")
	}

	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, _ := http.NewRequest("POST", "https://pkgstats.archlinux.de/post", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("pkgstats/%s", Version))
	response, err := c.Do(req)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	return string(body), nil
}
