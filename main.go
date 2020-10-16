package main

import (
	"flag"
	"fmt"
	"os"
)

// Version pkgstats version
var Version = "dev"

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
	packageChannel := async(pacman.GetInstalledPackages)
	mirrorChannel := async(pacman.GetServer)

	system := NewSystem()
	cpuArchitectureChannel := async(system.GetCpuArchitecture)
	architectureChannel := async(system.GetArchitecture)

	packages := <-packageChannel
	mirror := <-mirrorChannel
	cpuArchitecture := <-cpuArchitectureChannel
	architecture := <-architectureChannel

	if !*dryRun {
		client := NewClient(ApiBaseUrl)
		response, err := client.SendRequest(
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

		fmt.Printf("arch=%s\n", architecture)
		fmt.Printf("cpuarch=%s\n", cpuArchitecture)
		fmt.Printf("pkgstatsver=%s\n", Version)
		fmt.Printf("mirror=%s\n", mirror)
		fmt.Printf("quiet=%t\n", *quiet)
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

func async(f func() (string, error)) chan string {
	c := make(chan string)
	go func() {
		v, _ := f()
		c <- v
	}()
	return c
}
