// +build amd64

#include "textflag.h"

// func findKeyIndexAVX2(keys *[16]byte, key byte) int
TEXT ·findKeyIndexAVX2(SB), NOSPLIT, $0-16
	// Load arguments
	MOVQ keys+0(FP), SI    		// SI = keys array base pointer
	MOVB key+8(FP), AL      	// AL = key byte to find

	// Broadcast target byte to all lanes of YMM0
	MOVB AL, X0
	VPBROADCASTB X0, Y0

	// Load 16 bytes from keys array (we only need 16, not 32)
	VMOVDQU (SI), Y1

	// Compare with target byte using VPCMPEQB
	VPCMPEQB Y0, Y1, Y2

	// Extract the lower 128 bits (16 bytes) from YMM2 to XMM2
	VEXTRACTI128 $0, Y2, X2

	// Get mask of matching bytes (16 bits for 16 bytes)
	VPMOVMSKB X2, AX

	// Return the result
	MOVQ AX, ret+16(FP)
	RET

// func findInsertPositionAVX2(keys *[16]byte, key byte) int
TEXT ·findInsertPositionAVX2(SB), NOSPLIT, $0-16
	// Load arguments
	MOVQ keys+0(FP), SI    	// SI = keys array base pointer
	MOVB key+8(FP), AL      // AL = key byte to find

	// Broadcast target byte to all lanes of YMM0
	MOVB AL, X0
	VPBROADCASTB X0, Y0

	// Load 16 bytes from keys array
	VMOVDQU (SI), Y1

	// Compare target with keys (unsigned comparison for greater than)
	VPCMPGTB Y0, Y1, Y2

	// Extract the lower 128 bits (16 bytes) from YMM2 to XMM2
	VEXTRACTI128 $0, Y2, X2

	// Get mask of bytes greater than target
	VPMOVMSKB X2, AX

	// Return the result
	MOVQ AX, ret+16(FP)
	RET
