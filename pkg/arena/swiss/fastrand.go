//go:build !go1.22

package swiss

import (
	_ "unsafe"
)

//go:linkname fastrand runtime.fastrand
func fastrand() uint32

// randIntN returns a random number in the interval [0, n).
func randIntN(n int) uint32 {
	return fastModN(fastrand(), uint32(n))
}
