#include "textflag.h"

// func cpuid() uint32
TEXT ·cpuid(SB), NOSPLIT, $0-4
	MOVL $1, AX
	CPUID
	MOVL DX, ret+0(FP)
	RET
