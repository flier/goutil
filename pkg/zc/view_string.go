//go:build go1.20

package zc

import (
	"unsafe"

	"github.com/flier/goutil/pkg/xunsafe"
)

// String converts this View into a string, given its source.
func (r View) String(src *byte) string {
	if r.Len() == 0 {
		return ""
	}
	return unsafe.String(xunsafe.Add(src, r.Start()), r.Len())
}
