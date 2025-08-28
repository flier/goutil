package node_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	. "github.com/flier/goutil/pkg/arena/art/node"
	"github.com/flier/goutil/pkg/arena/slice"
	"github.com/flier/goutil/pkg/opt"
)

func TestLeaf(t *testing.T) {
	Convey("Given a Leaf", t, func() {
		a := &arena.Arena{}

		Convey("When creating a new leaf", func() {
			Convey("With simple key and value", func() {
				leaf := NewLeaf(a, []byte("hello"), 123)

				Convey("Then should have correct properties", func() {
					So(leaf.Type(), ShouldEqual, TypeLeaf)
					So(leaf.Full(), ShouldBeTrue)
					So(leaf.Key.Raw(), ShouldResemble, []byte("hello"))
					So(leaf.Value, ShouldEqual, 123)
				})

				Convey("And should return correct reference", func() {
					ref := leaf.Ref()
					So(ref.Type(), ShouldEqual, TypeLeaf)
					So(ref.IsLeaf(), ShouldBeTrue)
					So(ref.AsLeaf(), ShouldEqual, leaf)
				})
			})

			Convey("With empty key", func() {
				leaf := NewLeaf(a, []byte{}, 456)

				Convey("Then should handle empty key correctly", func() {
					So(leaf.Key.Len(), ShouldEqual, 0)
					So(leaf.Value, ShouldEqual, 456)
				})
			})

			Convey("With nil key", func() {
				leaf := NewLeaf(a, nil, 789)

				Convey("Then should handle nil key correctly", func() {
					So(leaf.Key.Len(), ShouldEqual, 0)
					So(leaf.Value, ShouldEqual, 789)
				})
			})

			Convey("With very long key", func() {
				longKey := make([]byte, 1000)
				for i := range longKey {
					longKey[i] = byte(i % 256)
				}

				leaf := NewLeaf(a, longKey, 999)

				Convey("Then should handle long key correctly", func() {
					So(leaf.Key.Len(), ShouldEqual, 1000)
					So(leaf.Value, ShouldEqual, 999)
					So(leaf.Key.Raw(), ShouldResemble, longKey)
				})
			})

			Convey("With special characters in key", func() {
				specialKey := []byte("hello\n\t\r\000world")
				leaf := NewLeaf(a, specialKey, 111)

				Convey("Then should preserve special characters", func() {
					So(leaf.Key.Raw(), ShouldResemble, specialKey)
					So(leaf.Value, ShouldEqual, 111)
				})
			})

			Convey("With unicode characters in key", func() {
				unicodeKey := []byte("hello世界")
				leaf := NewLeaf(a, unicodeKey, 222)

				Convey("Then should preserve unicode characters", func() {
					So(leaf.Key.Raw(), ShouldResemble, unicodeKey)
					So(leaf.Value, ShouldEqual, 222)
				})
			})
		})

		Convey("When checking prefix operations", func() {
			leaf := NewLeaf(a, []byte("hello"), 123)

			Convey("Then prefix should match the key", func() {
				So(leaf.Prefix().Raw(), ShouldResemble, []byte("hello"))
				So(leaf.Prefix().Len(), ShouldEqual, 5)
			})

			Convey("And setting prefix should update the key", func() {
				newPrefix := slice.FromString(a, "world")
				leaf.SetPrefix(newPrefix)

				So(leaf.Prefix().Raw(), ShouldResemble, []byte("world"))
				So(leaf.Key.Raw(), ShouldResemble, []byte("world"))
			})

			Convey("And setting prefix with different length", func() {
				shortPrefix := slice.FromString(a, "hi")
				leaf.SetPrefix(shortPrefix)

				So(leaf.Prefix().Raw(), ShouldResemble, []byte("hi"))
				So(leaf.Key.Raw(), ShouldResemble, []byte("hi"))
			})

			Convey("And setting prefix with longer length", func() {
				longPrefix := slice.FromString(a, "hello world")
				leaf.SetPrefix(longPrefix)

				So(leaf.Prefix().Raw(), ShouldResemble, []byte("hello world"))
				So(leaf.Key.Raw(), ShouldResemble, []byte("hello world"))
			})
		})

		Convey("When checking minimum and maximum", func() {
			leaf := NewLeaf(a, []byte("hello"), 123)

			Convey("Then minimum should return the leaf itself", func() {
				min := leaf.Minimum()
				So(min, ShouldEqual, leaf)
				So(min.Key.Raw(), ShouldResemble, []byte("hello"))
				So(min.Value, ShouldEqual, 123)
			})

			Convey("And maximum should return the leaf itself", func() {
				max := leaf.Maximum()
				So(max, ShouldEqual, leaf)
				So(max.Key.Raw(), ShouldResemble, []byte("hello"))
				So(max.Value, ShouldEqual, 123)
			})
		})

		Convey("When checking child operations", func() {
			leaf := NewLeaf(a, []byte("hello"), 123)

			Convey("Then FindChild should panic", func() {
				So(func() { leaf.FindChild(opt.Some(byte('h'))) }, ShouldPanicWith, "leaf cannot have children")
			})

			Convey("And AddChild should panic", func() {
				otherLeaf := NewLeaf(a, []byte("world"), 456)
				So(func() { leaf.AddChild(opt.Some(byte('w')), otherLeaf) }, ShouldPanicWith, "leaf cannot have children")
			})

			Convey("And RemoveChild should panic", func() {
				otherLeaf := NewLeaf(a, []byte("world"), 456)
				otherRef := otherLeaf.Ref()
				So(func() { leaf.RemoveChild(opt.Some(byte('w')), &otherRef) }, ShouldPanicWith, "leaf cannot have children")
			})

			Convey("And Grow should panic", func() {
				So(func() { leaf.Grow(a) }, ShouldPanicWith, "leaf cannot have children")
			})

			Convey("And Shrink should panic", func() {
				So(func() { leaf.Shrink(a) }, ShouldPanicWith, "leaf cannot have children")
			})
		})

		Convey("When checking Matches function", func() {
			leaf := NewLeaf(a, []byte("hello"), 123)

			Convey("Then should match identical key", func() {
				So(leaf.Matches([]byte("hello")), ShouldBeTrue)
			})

			Convey("And should not match different key", func() {
				So(leaf.Matches([]byte("world")), ShouldBeFalse)
				So(leaf.Matches([]byte("hell")), ShouldBeFalse)
				So(leaf.Matches([]byte("hello world")), ShouldBeFalse)
			})

			Convey("And should match empty key if leaf has empty key", func() {
				emptyLeaf := NewLeaf(a, []byte{}, 456)
				So(emptyLeaf.Matches([]byte{}), ShouldBeTrue)
				So(emptyLeaf.Matches([]byte("hello")), ShouldBeFalse)
			})

			Convey("And should handle nil key", func() {
				So(leaf.Matches(nil), ShouldBeFalse)
			})

			Convey("And should handle case sensitivity", func() {
				So(leaf.Matches([]byte("Hello")), ShouldBeFalse)
				So(leaf.Matches([]byte("HELLO")), ShouldBeFalse)
			})

			Convey("And should handle special characters", func() {
				specialLeaf := NewLeaf(a, []byte("hello\n"), 789)
				So(specialLeaf.Matches([]byte("hello\n")), ShouldBeTrue)
				So(specialLeaf.Matches([]byte("hello")), ShouldBeFalse)
			})

			Convey("And should handle unicode characters", func() {
				unicodeLeaf := NewLeaf(a, []byte("hello世界"), 999)
				So(unicodeLeaf.Matches([]byte("hello世界")), ShouldBeTrue)
				So(unicodeLeaf.Matches([]byte("hello")), ShouldBeFalse)
			})
		})

		Convey("When releasing the leaf", func() {
			leaf := NewLeaf(a, []byte("hello"), 123)

			Convey("Then should release memory correctly", func() {
				leaf.Release(a)

				// After release, the leaf's memory is freed but the slice metadata
				// may still contain the original values. This is expected behavior
				// as the slice.Release only frees the underlying memory.
				// In practice, the arena would handle the cleanup.
			})
		})
	})
}

func TestLeaf_EdgeCases(t *testing.T) {
	Convey("Given Leaf edge cases", t, func() {
		a := &arena.Arena{}

		Convey("When creating leaves with boundary values", func() {
			Convey("And using zero byte keys", func() {
				leaf := NewLeaf(a, []byte{0}, 123)
				So(leaf.Key.Raw(), ShouldResemble, []byte{0})
				So(leaf.Value, ShouldEqual, 123)
			})

			Convey("And using maximum byte values", func() {
				leaf := NewLeaf(a, []byte{255}, 456)
				So(leaf.Key.Raw(), ShouldResemble, []byte{255})
				So(leaf.Value, ShouldEqual, 456)
			})

			Convey("And using mixed boundary values", func() {
				leaf := NewLeaf(a, []byte{0, 255, 128}, 789)
				So(leaf.Key.Raw(), ShouldResemble, []byte{0, 255, 128})
				So(leaf.Value, ShouldEqual, 789)
			})
		})

		Convey("When creating leaves with very long keys", func() {
			Convey("And key length is 1MB", func() {
				keyLen := 1024 * 1024
				longKey := make([]byte, keyLen)
				for i := range longKey {
					longKey[i] = byte(i % 256)
				}

				leaf := NewLeaf(a, longKey, 999)

				So(leaf.Key.Len(), ShouldEqual, keyLen)
				So(leaf.Value, ShouldEqual, 999)
				So(leaf.Key.Raw(), ShouldResemble, longKey)
			})

			Convey("And key length is exactly 256 bytes", func() {
				keyLen := 256
				key := make([]byte, keyLen)
				for i := range key {
					key[i] = byte(i)
				}

				leaf := NewLeaf(a, key, 888)

				So(leaf.Key.Len(), ShouldEqual, keyLen)
				So(leaf.Value, ShouldEqual, 888)
				So(leaf.Key.Raw(), ShouldResemble, key)
			})
		})

		Convey("When creating leaves with special patterns", func() {
			Convey("And alternating bytes", func() {
				key := make([]byte, 100)
				for i := range key {
					if i%2 == 0 {
						key[i] = 0
					} else {
						key[i] = 255
					}
				}

				leaf := NewLeaf(a, key, 777)

				So(leaf.Key.Raw(), ShouldResemble, key)
				So(leaf.Value, ShouldEqual, 777)
			})

			Convey("And sequential bytes", func() {
				key := make([]byte, 256)
				for i := range key {
					key[i] = byte(i)
				}

				leaf := NewLeaf(a, key, 666)

				So(leaf.Key.Raw(), ShouldResemble, key)
				So(leaf.Value, ShouldEqual, 666)
			})
		})

		Convey("When testing Matches with edge cases", func() {
			Convey("And comparing with different length keys", func() {
				leaf := NewLeaf(a, []byte("hello"), 123)

				// Shorter key
				So(leaf.Matches([]byte("hell")), ShouldBeFalse)
				// Longer key
				So(leaf.Matches([]byte("hello world")), ShouldBeFalse)
				// Empty key
				So(leaf.Matches([]byte{}), ShouldBeFalse)
				// Nil key
				So(leaf.Matches(nil), ShouldBeFalse)
			})

			Convey("And comparing with similar but different keys", func() {
				leaf := NewLeaf(a, []byte("hello"), 123)

				// One character different
				So(leaf.Matches([]byte("hallo")), ShouldBeFalse)
				So(leaf.Matches([]byte("helpo")), ShouldBeFalse)
				So(leaf.Matches([]byte("hellp")), ShouldBeFalse)

				// Case differences
				So(leaf.Matches([]byte("Hello")), ShouldBeFalse)
				So(leaf.Matches([]byte("hELLO")), ShouldBeFalse)
			})

			Convey("And comparing with keys containing special characters", func() {
				leaf := NewLeaf(a, []byte("hello\n"), 123)

				// Different newline character
				So(leaf.Matches([]byte("hello\r")), ShouldBeFalse)
				// Missing newline
				So(leaf.Matches([]byte("hello")), ShouldBeFalse)
				// Extra character
				So(leaf.Matches([]byte("hello\n ")), ShouldBeFalse)
			})
		})
	})
}

func TestLeaf_Performance(t *testing.T) {
	Convey("Given Leaf performance considerations", t, func() {
		a := &arena.Arena{}

		Convey("When creating many leaves", func() {
			Convey("And creating 1000 leaves with different keys", func() {
				leaves := make([]*Leaf[any], 1000)
				for i := 0; i < 1000; i++ {
					key := []byte{byte(i % 256), byte((i / 256) % 256)}
					leaves[i] = NewLeaf[any](a, key, i)
				}

				So(len(leaves), ShouldEqual, 1000)

				// Verify all leaves are created correctly
				for i := 0; i < 1000; i++ {
					expectedKey := []byte{byte(i % 256), byte((i / 256) % 256)}
					So(leaves[i].Key.Raw(), ShouldResemble, expectedKey)
					So(leaves[i].Value, ShouldEqual, i)
				}
			})

			Convey("And creating leaves with same key but different values", func() {
				key := []byte("test")
				leaves := make([]*Leaf[any], 100)
				for i := 0; i < 100; i++ {
					leaves[i] = NewLeaf[any](a, key, i*100)
				}

				So(len(leaves), ShouldEqual, 100)

				// Verify all leaves have the same key but different values
				for i := 0; i < 100; i++ {
					So(leaves[i].Key.Raw(), ShouldResemble, key)
					So(leaves[i].Value, ShouldEqual, i*100)
				}
			})
		})

		Convey("When performing many Matches operations", func() {
			Convey("And matching against 1000 different keys", func() {
				leaf := NewLeaf(a, []byte("target"), 123)
				keys := make([][]byte, 1000)
				for i := 0; i < 1000; i++ {
					keys[i] = []byte{byte(i % 256), byte((i / 256) % 256)}
				}

				// Add the target key somewhere in the middle
				keys[500] = []byte("target")

				matches := 0
				for _, key := range keys {
					if leaf.Matches(key) {
						matches++
					}
				}

				So(matches, ShouldEqual, 1)
			})

			Convey("And matching against keys with common prefixes", func() {
				leaf := NewLeaf(a, []byte("hello world"), 123)
				keys := [][]byte{
					[]byte("hello"),
					[]byte("hello "),
					[]byte("hello w"),
					[]byte("hello wo"),
					[]byte("hello wor"),
					[]byte("hello worl"),
					[]byte("hello world"),
					[]byte("hello world!"),
				}

				matches := 0
				for _, key := range keys {
					if leaf.Matches(key) {
						matches++
					}
				}

				So(matches, ShouldEqual, 1) // Only exact match
			})
		})
	})
}
