#include "textflag.h"

// func cpuid() uint32
TEXT ·cpuid(SB), NOSPLIT, $0-4
	MOVL $0x80000001, AX
	MOVL $0, CX
	CPUID
	MOVL DX, ret+0(FP)
	RET
