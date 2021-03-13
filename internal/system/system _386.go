package system

import (
	"golang.org/x/sys/cpu"
)

func (system *System) GetCpuArchitecture() (architecture string, err error) {
	architecture = "i686"

	isX86_64 := cpu.X86.HasSSE2
	isX86_64V2 := isX86_64 && cpu.X86.HasPOPCNT && cpu.X86.HasSSE3 && cpu.X86.HasSSE41 && cpu.X86.HasSSE42 && cpu.X86.HasSSSE3
	isX86_64V3 := isX86_64V2 && cpu.X86.HasAVX && cpu.X86.HasAVX2 && cpu.X86.HasBMI1 && cpu.X86.HasBMI2 && cpu.X86.HasFMA && cpu.X86.HasOSXSAVE
	isX86_64V4 := isX86_64V3 && cpu.X86.HasAVX512F && cpu.X86.HasAVX512BW && cpu.X86.HasAVX512CD && cpu.X86.HasAVX512DQ && cpu.X86.HasAVX512VL

	if isX86_64V4 {
		architecture = "x86_64_v4"
	} else if isX86_64V3 {
		architecture = "x86_64_v3"
	} else if isX86_64V2 {
		architecture = "x86_64_v2"
	} else if isX86_64 {
		architecture = "x86_64"
	}

	return
}
