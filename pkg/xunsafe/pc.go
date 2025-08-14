package xunsafe

// PC is a raw function pointer, which can be used to store captureless
// funcs.
//
// Suppose a func() is in rax. Go implements calling it by emitting the
// following code:
//
//	mov  rdx, rax
//	mov  rcx, [rdx]
//	call rcx
//
// For a captureless func, this load will be of a constant containing the PC
// of the function to call. This can result in cache misses. This type works
// around that by keeping the PC local, so the resulting load avoids this
// problem.
type PC[F any] uintptr

// NewPC wraps a func. This performs no checking that the func does not
// capture any variables.
func NewPC[F any](f F) PC[F] {
	// Recall that a func()'s layout is *runtime.funcval, and PC[F] is emulating
	// runtime.funcval.
	return *BitCast[*PC[F]](f)
}

// Get returns the func this PC wraps.
func (pc *PC[F]) Get() F {
	return BitCast[F](pc)
}
