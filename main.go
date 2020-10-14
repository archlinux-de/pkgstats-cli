package main

import (
	"flag"
	"fmt"
	"os"
)

// Version pkgstats version
var Version = "2.5.0-0-dev"

// ApiBaseUrl pkgstats server URL
var ApiBaseUrl = "https://pkgstats.archlinux.de"

func main() {
	version := flag.Bool("v", false, "show the version of pkgstats")
	debug := flag.Bool("d", false, "enable debug mode")
	dryRun := flag.Bool("s", false, "show what information would be sent\n(but do not send anything)")
	quiet := flag.Bool("q", false, "be quiet except on errors")
	flag.Usage = printUsage
	flag.Parse()

	if os.Getenv("PKGSTATS_URL") != "" {
		ApiBaseUrl = os.Getenv("PKGSTATS_URL")
	}

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

	pacman := NewPacman()
	packages, _ := pacman.GetInstalledPackages()
	mirror, _ := pacman.GetServer()

	system := NewSystem()
	cpuArchitecture, _ := system.GetCpuArchitecture()
	architecture, _ := system.GetArchitecture()

	if !*dryRun {
		client := NewClient(ApiBaseUrl)
		response, err := client.sendRequest(
			packages,
			cpuArchitecture,
			architecture,
			mirror,
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
		fmt.Println(packages)

		fmt.Println("")

		fmt.Println("arch=", architecture)
		fmt.Println("cpuarch=", cpuArchitecture)
		fmt.Println("pkgstatsver=", Version)
		fmt.Println("mirror=", mirror)
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
	fmt.Printf("Statistics are available at %s\n", ApiBaseUrl)
}
