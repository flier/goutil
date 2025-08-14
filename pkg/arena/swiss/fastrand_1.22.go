//go:build go1.22

package swiss

import (
	"math/rand/v2"
)

// randIntN returns a random number in the interval [0, n).
func randIntN(n int) uint32 {
	return rand.Uint32N(uint32(n))
}
