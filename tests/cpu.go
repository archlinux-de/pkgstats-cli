package main

import (
	"fmt"
	"pkgstats-cli/internal/system"
)

func main() {
	system := system.NewSystem()
	cpu, _ := system.GetCpuArchitecture()
	fmt.Println(cpu)
}
