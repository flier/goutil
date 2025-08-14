// Copyright 2021-2025 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

//go:build !routinex

package g

import (
	"reflect"
	"unsafe"
)

// getg returns the pointer to the current runtime.g.
//
//go:nosplit
func getg() unsafe.Pointer

// getgp returns the pointer to the current runtime.g.
//
//go:nosplit
//go:linkname getgp runtime.getgp
func getgp() unsafe.Pointer {
	return getg()
}

// getgt returns the type of runtime.g.
//
//go:nosplit
//go:linkname getgt runtime.getgt
func getgt() reflect.Type {
	return typeByString("runtime.g")
}
