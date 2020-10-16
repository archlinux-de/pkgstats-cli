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

func TestHasLongMode(t *testing.T) {
	file := createCpuinfoMock(`flags:lm`)
	defer os.Remove(file)

	system := System{}
	system.cpuInfo = file

	res := system.hasLongMode()

	if !res {
		t.Error(res)
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

func TestGetCpuArchitecture(t *testing.T) {
	file := createCpuinfoMock(`flags:lm`)
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

func TestGetCpuArchitectureOn32BitOS(t *testing.T) {
	file := createCpuinfoMock(`flags:lm`)
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
