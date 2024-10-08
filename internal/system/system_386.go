package system

import (
	"golang.org/x/sys/cpu"
)

func (s *System) GetCpuArchitecture() (string, error) {
	// Check if we are on i686; otherwiese we know it is i586 with MMX as Go does require MMX instructions
	isI686 := hasCMOV()
	// We need to check for LM (Long Mode) as there are CPUs that support SSE2 but not x86_64 (e.g. Core Duo)
	isX86_64 := isI686 && hasLM() && cpu.X86.HasSSE2
	isX86_64V2 := isX86_64 && cpu.X86.HasPOPCNT && cpu.X86.HasSSE3 && cpu.X86.HasSSE41 && cpu.X86.HasSSE42 && cpu.X86.HasSSSE3
	isX86_64V3 := isX86_64V2 && cpu.X86.HasAVX && cpu.X86.HasAVX2 && cpu.X86.HasBMI1 && cpu.X86.HasBMI2 && cpu.X86.HasFMA && cpu.X86.HasOSXSAVE
	isX86_64V4 := isX86_64V3 && cpu.X86.HasAVX512F && cpu.X86.HasAVX512BW && cpu.X86.HasAVX512CD && cpu.X86.HasAVX512DQ && cpu.X86.HasAVX512VL

	switch {
	case isX86_64V4:
		return X86_64_V4, nil
	case isX86_64V3:
		return X86_64_V3, nil
	case isX86_64V2:
		return X86_64_V2, nil
	case isX86_64:
		return X86_64, nil
	case isI686:
		return I686, nil
	default:
		return I586, nil
	}
}

// Implemented at system_386.s
func hasLM() bool
func hasCMOV() bool
