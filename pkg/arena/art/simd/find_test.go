package simd

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Unit tests for FindKeyIndex using GoConvey
func TestFindKeyIndex(t *testing.T) {
	Convey("Given FindKeyIndex function", t, func() {
		Convey("When searching in an empty array", func() {
			keys := &[16]byte{}

			result := FindKeyIndex(keys, 0, 42)

			So(result, ShouldEqual, -1)
		})

		Convey("When searching in a single element array", func() {
			keys := &[16]byte{42}

			Convey("And the key is found", func() {
				result := FindKeyIndex(keys, 1, 42)

				So(result, ShouldEqual, 0)
			})

			Convey("And the key is not found", func() {
				result := FindKeyIndex(keys, 1, 24)

				So(result, ShouldEqual, -1)
			})
		})

		Convey("When searching in a multiple element array", func() {
			keys := &[16]byte{1, 2, 3, 4, 5}

			Convey("And searching for the first element", func() {
				result := FindKeyIndex(keys, 5, 1)

				So(result, ShouldEqual, 0)
			})

			Convey("And searching for a middle element", func() {
				result := FindKeyIndex(keys, 5, 3)

				So(result, ShouldEqual, 2)
			})

			Convey("And searching for the last element", func() {
				result := FindKeyIndex(keys, 5, 5)

				So(result, ShouldEqual, 4)
			})

			Convey("And searching for a non-existent element", func() {
				result := FindKeyIndex(keys, 5, 6)

				So(result, ShouldEqual, -1)
			})
		})

		Convey("When searching in a full 16-byte array", func() {
			keys := &[16]byte{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30}

			Convey("And searching for the first element", func() {
				result := FindKeyIndex(keys, 16, 0)

				So(result, ShouldEqual, 0)
			})

			Convey("And searching for a middle element", func() {
				result := FindKeyIndex(keys, 16, 16)

				So(result, ShouldEqual, 8)
			})

			Convey("And searching for the last element", func() {
				result := FindKeyIndex(keys, 16, 30)

				So(result, ShouldEqual, 15)
			})

			Convey("And searching for a non-existent element", func() {
				result := FindKeyIndex(keys, 16, 31)

				So(result, ShouldEqual, -1)
			})
		})

		Convey("When searching in an array with duplicate elements", func() {
			keys := &[16]byte{1, 1, 2, 3, 4}

			result := FindKeyIndex(keys, 5, 1)

			So(result, ShouldEqual, 0) // Should return first occurrence
		})

		Convey("When searching in an array with all same elements", func() {
			keys := &[16]byte{42, 42, 42, 42, 42}

			Convey("And searching for the existing element", func() {
				result := FindKeyIndex(keys, 5, 42)

				So(result, ShouldEqual, 0)
			})

			Convey("And searching for a non-existent element", func() {
				result := FindKeyIndex(keys, 5, 24)

				So(result, ShouldEqual, -1)
			})
		})
	})
}

// Unit tests for FindInsertPosition using GoConvey
func TestFindInsertPosition(t *testing.T) {
	Convey("Given FindInsertPosition function", t, func() {
		Convey("When inserting into an empty array", func() {
			keys := &[16]byte{}

			result := FindInsertPosition(keys, 0, 42)

			So(result, ShouldEqual, 0)
		})

		Convey("When inserting into a single element array", func() {
			keys := &[16]byte{5}

			Convey("And inserting before the existing element", func() {
				result := FindInsertPosition(keys, 1, 3)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting after the existing element", func() {
				result := FindInsertPosition(keys, 1, 7)

				So(result, ShouldEqual, 1)
			})
		})

		Convey("When inserting into a multiple element array", func() {
			keys := &[16]byte{2, 4, 6, 8}

			Convey("And inserting at the beginning", func() {
				result := FindInsertPosition(keys, 4, 1)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting in the middle", func() {
				result := FindInsertPosition(keys, 4, 5)

				So(result, ShouldEqual, 2)
			})

			Convey("And inserting at the end", func() {
				result := FindInsertPosition(keys, 4, 9)

				So(result, ShouldEqual, 4)
			})
		})

		Convey("When inserting into a full 16-byte array", func() {
			keys := &[16]byte{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32}

			Convey("And inserting at the beginning", func() {
				result := FindInsertPosition(keys, 16, 1)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting in the middle", func() {
				result := FindInsertPosition(keys, 16, 15)

				So(result, ShouldEqual, 7)
			})

			Convey("And inserting at the end", func() {
				result := FindInsertPosition(keys, 16, 33)

				So(result, ShouldEqual, 16)
			})
		})

		Convey("When inserting into an array with duplicate elements", func() {
			keys := &[16]byte{2, 2, 4, 6}

			Convey("And inserting before the first duplicate", func() {
				result := FindInsertPosition(keys, 4, 1)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting between duplicates", func() {
				result := FindInsertPosition(keys, 4, 3)

				So(result, ShouldEqual, 2)
			})
		})

		Convey("When inserting into an array with all same elements", func() {
			keys := &[16]byte{5, 5, 5, 5}

			Convey("And inserting before all elements", func() {
				result := FindInsertPosition(keys, 4, 3)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting after all elements", func() {
				result := FindInsertPosition(keys, 4, 7)

				So(result, ShouldEqual, 4)
			})
		})
	})
}

// Unit tests for scalar implementations using GoConvey
func TestFindKeyIndexScalar(t *testing.T) {
	Convey("Given findKeyIndexScalar function", t, func() {
		Convey("When searching in an empty array", func() {
			keys := &[16]byte{}

			result := findKeyIndexScalar(keys, 0, 42)

			So(result, ShouldEqual, -1)
		})

		Convey("When searching in a single element array", func() {
			keys := &[16]byte{42}

			Convey("And the key is found", func() {
				result := findKeyIndexScalar(keys, 1, 42)

				So(result, ShouldEqual, 0)
			})

			Convey("And the key is not found", func() {
				result := findKeyIndexScalar(keys, 1, 24)

				So(result, ShouldEqual, -1)
			})
		})

		Convey("When searching in a multiple element array", func() {
			keys := &[16]byte{1, 2, 3, 4, 5}

			Convey("And searching for an existing element", func() {
				result := findKeyIndexScalar(keys, 5, 3)

				So(result, ShouldEqual, 2)
			})

			Convey("And searching for a non-existent element", func() {
				result := findKeyIndexScalar(keys, 5, 6)

				So(result, ShouldEqual, -1)
			})
		})
	})
}

func TestFindInsertPositionScalar(t *testing.T) {
	Convey("Given findInsertPositionScalar function", t, func() {
		Convey("When inserting into an empty array", func() {
			keys := &[16]byte{}

			result := findInsertPositionScalar(keys, 0, 42)

			So(result, ShouldEqual, 0)
		})

		Convey("When inserting into a single element array", func() {
			keys := &[16]byte{5}

			Convey("And inserting before the existing element", func() {
				result := findInsertPositionScalar(keys, 1, 3)

				So(result, ShouldEqual, 0)
			})

			Convey("And inserting after the existing element", func() {
				result := findInsertPositionScalar(keys, 1, 7)

				So(result, ShouldEqual, 1)
			})
		})

		Convey("When inserting into a multiple element array", func() {
			keys := &[16]byte{2, 4, 6, 8}

			Convey("And inserting in the middle", func() {
				result := findInsertPositionScalar(keys, 4, 5)

				So(result, ShouldEqual, 2)
			})

			Convey("And inserting at the end", func() {
				result := findInsertPositionScalar(keys, 4, 9)

				So(result, ShouldEqual, 4)
			})
		})
	})
}

// Edge case tests using GoConvey
func TestEdgeCases(t *testing.T) {
	Convey("Given edge cases", t, func() {
		Convey("When testing boundary conditions", func() {
			keys := &[16]byte{255, 254, 253, 252}

			Convey("And testing maximum byte value", func() {
				result := FindKeyIndex(keys, 4, 255)
				So(result, ShouldEqual, 0)
			})

			Convey("And testing minimum byte value", func() {
				result := FindKeyIndex(keys, 4, 0)
				So(result, ShouldEqual, -1)
			})

			Convey("And testing insert position with maximum values", func() {
				result := FindInsertPosition(keys, 4, 255)
				So(result, ShouldEqual, 4)
			})

			Convey("And testing insert position with minimum values", func() {
				result := FindInsertPosition(keys, 4, 0)
				So(result, ShouldEqual, 0)
			})
		})

		Convey("When testing array bounds", func() {
			keys := &[16]byte{1, 2, 3, 4, 5}

			Convey("And testing with n = 0", func() {
				result := FindKeyIndex(keys, 0, 1)
				So(result, ShouldEqual, -1)

				result = FindInsertPosition(keys, 0, 1)
				So(result, ShouldEqual, 0)
			})

			Convey("And testing with n = 1", func() {
				result := FindKeyIndex(keys, 1, 1)
				So(result, ShouldEqual, 0)

				result = FindInsertPosition(keys, 1, 0)
				So(result, ShouldEqual, 0)
			})
		})
	})
}

// Test for FindNonZeroKeyIndex function
func TestFindNonZeroKeyIndex(t *testing.T) {
	Convey("Given FindNonZeroKeyIndex function", t, func() {
		Convey("When searching in an array with all zeros", func() {
			keys := &[256]byte{}

			result := FindNonZeroKeyIndex(keys)

			So(result, ShouldEqual, -1)
		})

		Convey("When searching in an array with first element non-zero", func() {
			keys := &[256]byte{}
			keys[0] = 42

			result := FindNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 0)
		})

		Convey("When searching in an array with middle element non-zero", func() {
			keys := &[256]byte{}
			keys[128] = 42

			result := FindNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 128)
		})

		Convey("When searching in an array with last element non-zero", func() {
			keys := &[256]byte{}
			keys[255] = 42

			result := FindNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 255)
		})

		Convey("When searching in an array with multiple non-zero elements", func() {
			keys := &[256]byte{}
			keys[10] = 42
			keys[20] = 100
			keys[30] = 200

			result := FindNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 10) // Should return first non-zero
		})
	})
}

// Test for FindLastNonZeroKeyIndex function
func TestFindLastNonZeroKeyIndex(t *testing.T) {
	Convey("Given FindLastNonZeroKeyIndex function", t, func() {
		Convey("When searching in an array with all zeros", func() {
			keys := &[256]byte{}

			result := FindLastNonZeroKeyIndex(keys)

			So(result, ShouldEqual, -1)
		})

		Convey("When searching in an array with first element non-zero", func() {
			keys := &[256]byte{}
			keys[0] = 42

			result := FindLastNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 0) // Should return 0 since it's the last (and only) non-zero element
		})

		Convey("When searching in an array with middle element non-zero", func() {
			keys := &[256]byte{}
			keys[128] = 42

			result := FindLastNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 128)
		})

		Convey("When searching in an array with last element non-zero", func() {
			keys := &[256]byte{}
			keys[255] = 42

			result := FindLastNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 255)
		})

		Convey("When searching in an array with multiple non-zero elements", func() {
			keys := &[256]byte{}
			keys[10] = 42
			keys[20] = 100
			keys[30] = 200

			result := FindLastNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 30) // Should return last non-zero
		})

		Convey("When searching in an array with non-zero elements at boundaries", func() {
			keys := &[256]byte{}
			keys[0] = 1
			keys[255] = 255

			result := FindLastNonZeroKeyIndex(keys)

			So(result, ShouldEqual, 255) // Should return last non-zero
		})
	})
}
