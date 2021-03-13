package system

import (
	. "golang.org/x/sys/cpu"
)

func (system *System) GetCpuArchitecture() (architecture string, err error) {
	architecture = "i686"

	// https://gitlab.com/x86-psABIs/x86-64-ABI/-/blob/master/x86-64-ABI/low-level-sys-info.tex
	// We cannot use intel/cpuid on i686, so we fallback on a subset of flags that are provided by sys/cpu
	isX86_64 := X86.HasSSE2
	isX86_64V2 := isX86_64 && X86.HasPOPCNT && X86.HasSSE3 && X86.HasSSE41 && X86.HasSSE42 && X86.HasSSSE3
	isX86_64V3 := isX86_64V2 && X86.HasAVX && X86.HasAVX2 && X86.HasBMI1 && X86.HasBMI2 && X86.HasFMA && X86.HasOSXSAVE
	isX86_64V4 := isX86_64V3 && X86.HasAVX512F && X86.HasAVX512BW && X86.HasAVX512CD && X86.HasAVX512DQ && X86.HasAVX512VL

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
