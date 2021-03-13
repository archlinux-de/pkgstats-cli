package system

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func init() {
	Mocks["uname"] = func() {
		fmt.Println("x86_64")
		os.Exit(0)
	}

	Mocks["uname32"] = func() {
		fmt.Println("i686")
		os.Exit(0)
	}
}

var Mocks = make(map[string]func())

func TestMain(m *testing.M) {
	mockName := os.Getenv("TEST_MOCK")
	if mockName != "" {
		mock, ok := Mocks[mockName]
		if ok {
			mock()
		}
	}

	os.Exit(m.Run())
}
func TestGetMachine(t *testing.T) {
	system := System{}
	system.uname = os.Args[0]
	system.env = []string{"TEST_MOCK=uname"}

	out, err := system.getMachine()

	if err != nil {
		t.Error(err)
	}
	if out != "x86_64" {
		t.Error(out)
	}
}

func createCpuinfoMock(cpuinfo string) string {
	file, err := ioutil.TempFile("", "cpuinfo")
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(cpuinfo)
	if err != nil {
		log.Fatal(err)
	}

	return file.Name()
}

func TestGetCpuArchitectureOn32BitOS(t *testing.T) {
	file := createCpuinfoMock(`flags:lm cmov cx8 fpu fxsr mmx syscall sse sse2`)
	defer os.Remove(file)

	system := System{}
	system.uname = os.Args[0]
	system.cpuInfo = file
	system.env = []string{"TEST_MOCK=uname32"}

	out, err := system.GetCpuArchitecture()

	if err != nil {
		t.Error(err)
	}
	if out != "x86_64" {
		t.Error(out)
	}
}

func TestGetCpuArchitectureOn32BitCpu(t *testing.T) {
	file := createCpuinfoMock(`flags:mmx`)
	defer os.Remove(file)

	system := System{}
	system.uname = os.Args[0]
	system.cpuInfo = file
	system.env = []string{"TEST_MOCK=uname32"}

	out, err := system.GetCpuArchitecture()

	if err != nil {
		t.Error(err)
	}
	if out != "i686" {
		t.Error(out, file)
	}
}

func TestGetCpuArchitectureX86_64(t *testing.T) {
	file := createCpuinfoMock(`flags:lm cmov cx8 fpu fxsr mmx syscall sse sse2`)
	defer os.Remove(file)

	system := System{}
	system.uname = os.Args[0]
	system.cpuInfo = file
	system.env = []string{"TEST_MOCK=uname"}

	out, err := system.GetCpuArchitecture()

	if err != nil {
		t.Error(err)
	}
	if out != "x86_64" {
		t.Error(out)
	}
}

func TestGetCpuArchitectureX86_64V2(t *testing.T) {
	file := createCpuinfoMock(`flags:lm cmov cx8 fpu fxsr mmx syscall sse sse2 cx16 lahf_lm popcnt pni sse4_1 sse4_2 ssse3`)
	defer os.Remove(file)

	system := System{}
	system.uname = os.Args[0]
	system.cpuInfo = file
	system.env = []string{"TEST_MOCK=uname"}

	out, err := system.GetCpuArchitecture()

	if err != nil {
		t.Error(err)
	}
	if out != "x86_64_v2" {
		t.Error(out)
	}
}

func TestGetCpuArchitectureX86_64V3(t *testing.T) {
	file := createCpuinfoMock(`flags:lm cmov cx8 fpu fxsr mmx syscall sse sse2 cx16 lahf_lm popcnt pni sse4_1 sse4_2 ssse3 avx avx2 bmi1 bmi2 f16c fma abm movbe xsave`)
	defer os.Remove(file)

	system := System{}
	system.uname = os.Args[0]
	system.cpuInfo = file
	system.env = []string{"TEST_MOCK=uname"}

	out, err := system.GetCpuArchitecture()

	if err != nil {
		t.Error(err)
	}
	if out != "x86_64_v3" {
		t.Error(out)
	}
}

func TestGetCpuArchitectureX86_64V4(t *testing.T) {
	file := createCpuinfoMock(`flags:lm cmov cx8 fpu fxsr mmx syscall sse sse2 cx16 lahf_lm popcnt pni sse4_1 sse4_2 ssse3 avx avx2 bmi1 bmi2 f16c fma abm movbe xsave avx512f avx512bw avx512cd avx512dq avx512vl`)
	defer os.Remove(file)

	system := System{}
	system.uname = os.Args[0]
	system.cpuInfo = file
	system.env = []string{"TEST_MOCK=uname"}

	out, err := system.GetCpuArchitecture()

	if err != nil {
		t.Error(err)
	}
	if out != "x86_64_v4" {
		t.Error(out)
	}
}
