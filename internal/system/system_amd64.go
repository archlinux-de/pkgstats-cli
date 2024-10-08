package system

import (
	"golang.org/x/sys/cpu"
)

func (s *System) GetCpuArchitecture() (string, error) {
	// https://gitlab.com/x86-psABIs/x86-64-ABI/-/blob/master/x86-64-ABI/low-level-sys-info.tex
	// https://unix.stackexchange.com/questions/43539/what-do-the-flags-in-proc-cpuinfo-mean/43540#43540
	isX86_64V2 := cpu.X86.HasPOPCNT && cpu.X86.HasSSE3 && cpu.X86.HasSSE41 && cpu.X86.HasSSE42 && cpu.X86.HasSSSE3
	isX86_64V3 := isX86_64V2 && cpu.X86.HasAVX && cpu.X86.HasAVX2 && cpu.X86.HasBMI1 && cpu.X86.HasBMI2 && cpu.X86.HasFMA && cpu.X86.HasOSXSAVE
	isX86_64V4 := isX86_64V3 && cpu.X86.HasAVX512F && cpu.X86.HasAVX512BW && cpu.X86.HasAVX512CD && cpu.X86.HasAVX512DQ && cpu.X86.HasAVX512VL

	switch {
	case isX86_64V4:
		return X86_64_V4, nil
	case isX86_64V3:
		return X86_64_V3, nil
	case isX86_64V2:
		return X86_64_V2, nil
	default:
		return X86_64, nil
	}
}
