package system

import (
	"errors"

	"github.com/intel-go/cpuid"
)

func (system *System) GetCpuArchitecture() (architecture string, err error) {
	// https://gitlab.com/x86-psABIs/x86-64-ABI/-/blob/master/x86-64-ABI/low-level-sys-info.tex
	// https://unix.stackexchange.com/questions/43539/what-do-the-flags-in-proc-cpuinfo-mean/43540#43540
	isX86_64 := system.hasFeatures([]uint64{cpuid.LM, cpuid.CMOV, cpuid.CX8, cpuid.FPU, cpuid.FXSR, cpuid.MMX, cpuid.SYSCALL, cpuid.SSE, cpuid.SSE2})
	isX86_64V2 := isX86_64 && system.hasFeatures([]uint64{cpuid.CX16, cpuid.LAHF_LM, cpuid.POPCNT, cpuid.SSE3, cpuid.SSE4_1, cpuid.SSE4_2, cpuid.SSSE3})
	isX86_64V3 := isX86_64V2 && system.hasFeatures([]uint64{cpuid.AVX, cpuid.AVX2, cpuid.BMI1, cpuid.BMI2, cpuid.F16C, cpuid.FMA, cpuid.ABM, cpuid.MOVBE, cpuid.OSXSAVE})
	isX86_64V4 := isX86_64V3 && system.hasFeatures([]uint64{cpuid.AVX512F, cpuid.AVX512BW, cpuid.AVX512CD, cpuid.AVX512DQ, cpuid.AVX512VL})

	if isX86_64V4 {
		architecture = "x86_64_v4"
	} else if isX86_64V3 {
		architecture = "x86_64_v3"
	} else if isX86_64V2 {
		architecture = "x86_64_v2"
	} else if isX86_64 {
		architecture = "x86_64"
	} else {
		err = errors.New("unknown CPU architecture")
	}
	return
}

func (system *System) hasFeatures(features []uint64) bool {
	for _, feature := range features {
		if !cpuid.HasFeature(feature) && !cpuid.HasExtendedFeature(feature) && !cpuid.HasExtraFeature(feature) {
			return false
		}
	}
	return true
}
