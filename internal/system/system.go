package system

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strings"
)

type System struct {
	env     []string
	uname   string
	cpuInfo string
}

func NewSystem() System {
	system := System{}
	system.uname = "uname"
	system.cpuInfo = "/proc/cpuinfo"
	return system
}

func (system *System) GetArchitecture() (string, error) {
	arch, err := system.getMachine()
	return arch, err
}

func (system *System) GetCpuArchitecture() (string, error) {
	architecture, err := system.GetArchitecture()

	cpuFlags, err := system.getCPUFlags()
	// detect a 64 bit CPU when ruinning a 32 bit OS
	if architecture == "i686" && system.inArray("lm", cpuFlags) {
		architecture = "x86_64"
	}

	if architecture == "x86_64" {
		// detect different levels of x86_64
		// https://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/tree/arch/x86/include/asm/cpufeatures.h
		// https://gitlab.com/x86-psABIs/x86-64-ABI/-/blob/master/x86-64-ABI/low-level-sys-info.tex
		// https://unix.stackexchange.com/questions/43539/what-do-the-flags-in-proc-cpuinfo-mean/43540#43540
		isx86_64 := system.inArray("lm", cpuFlags)
		isx86_64_V1 := isx86_64 && system.arrayInArray([]string{"cmov", "cx8", "fpu", "fxsr", "mmx", "syscall", "sse", "sse2"}, cpuFlags)
		isx86_64_V2 := isx86_64_V1 && system.arrayInArray([]string{"cx16", "lahf_lm", "popcnt", "pni", "sse4_1", "sse4_2", "ssse3"}, cpuFlags)
		isx86_64_V3 := isx86_64_V2 && system.arrayInArray([]string{"avx", "avx2", "bmi1", "bmi2", "f16c", "fma", "abm", "movbe", "xsave"}, cpuFlags)
		isx86_64_V4 := isx86_64_V3 && system.arrayInArray([]string{"avx512f", "avx512bw", "avx512cd", "avx512dq", "avx512vl"}, cpuFlags)

		if isx86_64_V4 {
			architecture = "x86_64_v4"
		} else if isx86_64_V3 {
			architecture = "x86_64_v3"
		} else if isx86_64_V2 {
			architecture = "x86_64_v2"
		} else if isx86_64_V1 {
			architecture = "x86_64"
		}
	}

	return architecture, err
}

func (system *System) getMachine() (string, error) {
	cmd := exec.Command(system.uname, "-m")
	cmd.Env = system.env
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func (system *System) getCPUFlags() (info []string, err error) {
	cpuInfo, err := os.Open(system.cpuInfo)
	if err != nil {
		return []string{}, err
	}
	defer cpuInfo.Close()

	scanner := bufio.NewScanner(cpuInfo)
	for scanner.Scan() {
		newline := scanner.Text()
		list := strings.Split(newline, ":")

		if len(list) > 1 && strings.EqualFold(strings.TrimSpace(list[0]), "flags") {
			return strings.Fields(list[1]), nil
		}
	}

	err = scanner.Err()
	if err != nil {
		return []string{}, err
	}

	return []string{}, errors.New("No CPU flags found")
}

func (system *System) inArray(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func (system *System) arrayInArray(needles []string, haystack []string) bool {
	for _, needle := range needles {
		if !system.inArray(needle, haystack) {
			return false
		}
	}
	return true
}
