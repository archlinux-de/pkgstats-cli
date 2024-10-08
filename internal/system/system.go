package system

type System struct{}

const (
	AARCH64   = "aarch64"
	ARMV5     = "armv5"
	ARMV6     = "armv6"
	ARMV7     = "armv7"
	I586      = "i586"
	I686      = "i686"
	X86_64    = "x86_64"
	X86_64_V2 = "x86_64_v2"
	X86_64_V3 = "x86_64_v3"
	X86_64_V4 = "x86_64_v4"
)

func NewSystem() *System {
	return &System{}
}
