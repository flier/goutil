package xunsafe_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/xunsafe"
)

func TestEscape(t *testing.T) {
	Convey("Given escape analysis operations", t, func() {
		Convey("When escaping pointers", func() {
			Convey("And escaping int pointer", func() {
				i := 42
				escaped := xunsafe.Escape(&i)
				So(escaped, ShouldEqual, &i)
			})

			Convey("And escaping string pointer", func() {
				s := "hello"
				escapedStr := xunsafe.Escape(&s)
				So(escapedStr, ShouldEqual, &s)
			})

			Convey("And escaping struct pointer", func() {
				type TestStruct struct {
					ID   int
					Name string
				}
				ts := TestStruct{ID: 1, Name: "test"}
				escapedStruct := xunsafe.Escape(&ts)
				So(escapedStruct, ShouldEqual, &ts)
			})

			Convey("And escaping nil pointer", func() {
				var nilPtr *int
				escapedNil := xunsafe.Escape(nilPtr)
				So(escapedNil, ShouldEqual, nilPtr)
			})
		})

		Convey("When preventing escape", func() {
			Convey("And preventing escape of int pointer", func() {
				i := 42
				noEscaped := xunsafe.NoEscape(&i)
				So(noEscaped, ShouldEqual, &i)
			})

			Convey("And preventing escape of string pointer", func() {
				s := "hello"
				noEscapedStr := xunsafe.NoEscape(&s)
				So(noEscapedStr, ShouldEqual, &s)
			})

			Convey("And preventing escape of struct pointer", func() {
				type TestStruct struct {
					ID   int
					Name string
				}
				ts := TestStruct{ID: 1, Name: "test"}
				noEscapedStruct := xunsafe.NoEscape(&ts)
				So(noEscapedStruct, ShouldEqual, &ts)
			})

			Convey("And preventing escape of nil pointer", func() {
				var nilPtr *int
				noEscapedNil := xunsafe.NoEscape(nilPtr)
				So(noEscapedNil, ShouldEqual, nilPtr)
			})
		})

		Convey("When handling edge cases", func() {
			Convey("And escaping zero value", func() {
				var zeroInt int
				escapedZero := xunsafe.Escape(&zeroInt)
				So(escapedZero, ShouldEqual, &zeroInt)
			})

			Convey("And escaping empty string", func() {
				emptyStr := ""
				escapedEmpty := xunsafe.Escape(&emptyStr)
				So(escapedEmpty, ShouldEqual, &emptyStr)
			})

			Convey("And escaping empty struct", func() {
				type EmptyStruct struct{}
				empty := EmptyStruct{}
				escapedEmptyStruct := xunsafe.Escape(&empty)
				So(escapedEmptyStruct, ShouldEqual, &empty)
			})

			Convey("And preventing escape of zero value", func() {
				var zeroInt int
				noEscapedZero := xunsafe.NoEscape(&zeroInt)
				So(noEscapedZero, ShouldEqual, &zeroInt)
			})

			Convey("And preventing escape of empty string", func() {
				emptyStr := ""
				noEscapedEmpty := xunsafe.NoEscape(&emptyStr)
				So(noEscapedEmpty, ShouldEqual, &emptyStr)
			})

			Convey("And preventing escape of empty struct", func() {
				type EmptyStruct struct{}
				empty := EmptyStruct{}
				noEscapedEmptyStruct := xunsafe.NoEscape(&empty)
				So(noEscapedEmptyStruct, ShouldEqual, &empty)
			})
		})

		Convey("When performing comprehensive tests", func() {
			Convey("Given various types", func() {
				testCases := []struct {
					name  string
					value any
				}{
					{"int", 42},
					{"string", "hello"},
					{"slice", []int{1, 2, 3}},
					{"map", map[string]int{"a": 1}},
					{"struct", struct{ ID int }{ID: 42}},
					{"pointer", &struct{}{}},
				}

				for _, tc := range testCases {
					Convey(fmt.Sprintf("When testing %s", tc.name), func() {
						Convey("And testing Escape", func() {
							escaped := xunsafe.Escape(&tc.value)
							So(escaped, ShouldEqual, &tc.value)
						})

						Convey("And testing NoEscape", func() {
							noEscaped := xunsafe.NoEscape(&tc.value)
							So(noEscaped, ShouldEqual, &tc.value)
						})
					})
				}
			})
		})
	})
}
