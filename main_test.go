package main

import "testing"

func Test_GetArchitecture(t *testing.T) {
	architecture, _ := getArchitecture()
	if architecture != "x86_64" {
		t.Error("Unexptected architecture")
	}
}

func Test_GetCpuArchitecture(t *testing.T) {
	cpuArchitecture, _ := getCpuArchitecture()
	if cpuArchitecture != "x86_64" {
		t.Error("Unexptected cpu architecture")
	}
}

func Test_GetMirror(t *testing.T) {
	mirror, _ := getMirror()
	if mirror == "" {
		t.Error("No mirror found")
	}
}

func Test_GetPackages(t *testing.T) {
	packages, _ := getPackages()
	if packages == "" {
		t.Error("No packages found")
	}
}
