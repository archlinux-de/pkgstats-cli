package system

import (
	"golang.org/x/sys/cpu"
)

func (system *System) GetCpuArchitecture() (architecture string, err error) {
	architecture = "armv5"

	// https://github.com/lpereira/hardinfo/blob/master/modules/devices/arm/processor.c#L180
	// https://github.com/golang/go/issues/38987#issuecomment-626513091
	if cpu.ARM.HasPMULL || cpu.ARM.HasCRC32 {
		architecture = "aarch64"
	} else if cpu.ARM.HasVFPv3 {
		architecture = "armv7"
	} else if cpu.ARM.HasVFP {
		architecture = "armv6"
	}

	return
}
