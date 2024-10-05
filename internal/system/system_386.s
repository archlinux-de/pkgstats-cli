#include "textflag.h"

// func hasLM() bool
TEXT ·hasLM(SB), NOSPLIT, $0-1
    MOVL $0x80000001, AX  // Use extended function 0x80000001
    CPUID
    SHRL $29, DX          // Shift right by 29 bits
    ANDL $1, DX           // Isolate the 29th bit
    MOVB DX, ret+0(FP)    // Move the result to the return value
    RET

// func hasCMOV() bool
TEXT ·hasCMOV(SB), NOSPLIT, $0-1
    MOVL $1, AX  // Use standard function 1
    CPUID
    SHRL $15, DX // Shift right by 15 bits
    ANDL $1, DX  // Isolate the 15th bit
    MOVB DX, ret+0(FP)  // Move the result to the return value
    RET
