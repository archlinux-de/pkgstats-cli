package system

import (
	"errors"

	. "github.com/intel-go/cpuid"
)

func (system *System) GetCpuArchitecture() (architecture string, err error) {
	// https://gitlab.com/x86-psABIs/x86-64-ABI/-/blob/master/x86-64-ABI/low-level-sys-info.tex
	// https://unix.stackexchange.com/questions/43539/what-do-the-flags-in-proc-cpuinfo-mean/43540#43540
	isX86_64 := system.hasFeatures([]uint64{LM, CMOV, CX8, FPU, FXSR, MMX, SYSCALL, SSE, SSE2})
	isX86_64V2 := isX86_64 && system.hasFeatures([]uint64{CX16, LAHF_LM, POPCNT, SSE3, SSE4_1, SSE4_2, SSSE3})
	isX86_64V3 := isX86_64V2 && system.hasFeatures([]uint64{AVX, AVX2, BMI1, BMI2, F16C, FMA, ABM, MOVBE, OSXSAVE})
	isX86_64V4 := isX86_64V3 && system.hasFeatures([]uint64{AVX512F, AVX512BW, AVX512CD, AVX512DQ, AVX512VL})

	if isX86_64V4 {
		architecture = "x86_64_v4"
	} else if isX86_64V3 {
		architecture = "x86_64_v3"
	} else if isX86_64V2 {
		architecture = "x86_64_v2"
	} else if isX86_64 {
		architecture = "x86_64"
	} else {
		err = errors.New("Unknown CPU architecture")
	}
	return
}

func (system *System) hasFeatures(features []uint64) bool {
	for _, feature := range features {
		if !HasFeature(feature) && !HasExtendedFeature(feature) && !HasExtraFeature(feature) {
			return false
		}
	}
	return true
}
