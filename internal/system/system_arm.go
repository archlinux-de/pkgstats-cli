package system

import (
	. "golang.org/x/sys/cpu"
)

func (s *System) GetCpuArchitecture() (string, error) {
	// https://github.com/lpereira/hardinfo/blob/master/modules/devices/arm/processor.c#L180
	// https://github.com/golang/go/issues/38987#issuecomment-626513091
	// https://community.arm.com/developer/tools-software/oss-platforms/b/android-blog/posts/runtime-detection-of-cpu-features-on-an-armv8-a-cpu
	// https://developer.arm.com/documentation/dui0471/m/key-features-of-arm-architecture-versions/arm-architecture-v7-a?lang=en
	switch {
	case ARM.HasPMULL || ARM.HasCRC32 || ARM.HasAES || ARM.HasSHA1 || ARM.HasSHA2:
		return AARCH64, nil
	case ARM.HasVFPv3 && ARM.HasTHUMBEE:
		return ARMV7, nil
	case ARM.HasVFP && ARM.HasTHUMB:
		return ARMV6, nil
	default:
		return ARMV5, nil
	}
}
