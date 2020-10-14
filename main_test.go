package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestMain(t *testing.M) {
	_, err := exec.LookPath("pacman")
	if err == nil {
		t.Run()
	} else {
		fmt.Println("Testing is only supported on Arch Linux", err)
		os.Exit(0)
	}
}

func Test_GetArchitecture(t *testing.T) {
	architecture, err := getArchitecture()
	if architecture != "x86_64" || err != nil {
		t.Error("Unexptected architecture", err)
	}
}

func Test_GetCpuArchitecture(t *testing.T) {
	cpuArchitecture, err := getCpuArchitecture()
	if cpuArchitecture != "x86_64" || err != nil {
		t.Error("Unexptected cpu architecture", err)
	}
}

func Test_GetMirror(t *testing.T) {
	mirror, err := getMirror()
	if mirror == "" || err != nil {
		t.Error("No mirror found", err)
	}
}

func Test_GetPackages(t *testing.T) {
	packages, err := getPackages()
	if packages == "" || err != nil {
		t.Error("No packages found", err)
	}
}
