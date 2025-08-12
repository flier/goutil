//go:build go1.23

package xiter_test

import (
	"fmt"
	"slices"
	"testing"
	"unicode/utf8"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/flier/goutil/pkg/xiter"
)

func ExampleChars() {
	fmt.Println(slices.Collect(Chars([]byte("héllo世界"))))
	// Output:
	// [104 233 108 108 111 19990 30028]
}

func TestChars(t *testing.T) {
	Convey("Chars", t, func() {
		Convey("should iterate over runes", func() {
			input := []byte("héllo世界")
			want := []rune{'h', 'é', 'l', 'l', 'o', '世', '界'}

			So(slices.Collect(Chars(input)), ShouldResemble, want)
		})

		Convey("should handle empty", func() {

			So(slices.Collect(Chars([]byte(""))), ShouldBeNil)
		})

		Convey("should handle invalid UTF-8", func() {
			input := []byte{0xff, 'a'}
			want := []rune{utf8.RuneError, 'a'}

			So(slices.Collect(Chars(input)), ShouldResemble, want)
		})
	})
}
