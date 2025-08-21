package node_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	. "github.com/flier/goutil/pkg/arena/art/node"
)

func TestRef(t *testing.T) {
	Convey("Given Ref type", t, func() {
		a := &arena.Arena{}

		Convey("When creating references", func() {
			Convey("And creating leaf reference", func() {
				leaf := NewLeaf[any](a, []byte("hello"), 123)
				So(leaf, ShouldNotBeNil)

				ref := NewRef[any](TypeLeaf, leaf)

				Convey("Then should have correct properties", func() {
					So(ref.Type(), ShouldEqual, TypeLeaf)
					So(ref.IsLeaf(), ShouldBeTrue)
					So(ref.IsNode(), ShouldBeFalse)
					So(ref.Empty(), ShouldBeFalse)
				})

				Convey("And should return correct leaf", func() {
					leafResult := ref.AsLeaf()
					So(leafResult, ShouldNotBeNil)
					So(leafResult.Key.Raw(), ShouldResemble, []byte("hello"))
					So(leafResult.Value, ShouldEqual, 123)
				})

				Convey("And should return nil for node types", func() {
					So(ref.AsNode4(), ShouldBeNil)
					So(ref.AsNode16(), ShouldBeNil)
					So(ref.AsNode48(), ShouldBeNil)
					So(ref.AsNode256(), ShouldBeNil)
				})
			})

			Convey("And creating Node4 reference", func() {
				node4 := arena.New(a, Node4[any]{})
				ref := NewRef[any](TypeNode4, node4)

				Convey("Then should have correct properties", func() {
					So(ref.Type(), ShouldEqual, TypeNode4)
					So(ref.IsNode4(), ShouldBeTrue)
					So(ref.IsNode(), ShouldBeTrue)
					So(ref.Empty(), ShouldBeFalse)
				})

				Convey("And should return correct node", func() {
					nodeResult := ref.AsNode4()
					So(nodeResult, ShouldNotBeNil)
					So(nodeResult, ShouldEqual, node4)
				})

				Convey("And should return nil for other node types", func() {
					So(ref.AsLeaf(), ShouldBeNil)
					So(ref.AsNode16(), ShouldBeNil)
					So(ref.AsNode48(), ShouldBeNil)
					So(ref.AsNode256(), ShouldBeNil)
				})
			})

			Convey("And creating Node16 reference", func() {
				node16 := arena.New(a, Node16[any]{})
				ref := NewRef[any](TypeNode16, node16)

				Convey("Then should have correct properties", func() {
					So(ref.Type(), ShouldEqual, TypeNode16)
					So(ref.IsNode16(), ShouldBeTrue)
					So(ref.IsNode(), ShouldBeTrue)
					So(ref.Empty(), ShouldBeFalse)
				})

				Convey("And should return correct node", func() {
					nodeResult := ref.AsNode16()
					So(nodeResult, ShouldNotBeNil)
					So(nodeResult, ShouldEqual, node16)
				})
			})

			Convey("And creating Node48 reference", func() {
				node48 := arena.New(a, Node48[any]{})
				ref := NewRef[any](TypeNode48, node48)

				Convey("Then should have correct properties", func() {
					So(ref.Type(), ShouldEqual, TypeNode48)
					So(ref.IsNode48(), ShouldBeTrue)
					So(ref.IsNode(), ShouldBeTrue)
					So(ref.Empty(), ShouldBeFalse)
				})

				Convey("And should return correct node", func() {
					nodeResult := ref.AsNode48()
					So(nodeResult, ShouldNotBeNil)
					So(nodeResult, ShouldEqual, node48)
				})
			})

			Convey("And creating Node256 reference", func() {
				node256 := arena.New(a, Node256[any]{})
				ref := NewRef[any](TypeNode256, node256)

				Convey("Then should have correct properties", func() {
					So(ref.Type(), ShouldEqual, TypeNode256)
					So(ref.IsNode256(), ShouldBeTrue)
					So(ref.IsNode(), ShouldBeTrue)
					So(ref.Empty(), ShouldBeFalse)
				})

				Convey("And should return correct node", func() {
					nodeResult := ref.AsNode256()
					So(nodeResult, ShouldNotBeNil)
					So(nodeResult, ShouldEqual, node256)
				})
			})
		})

		Convey("When working with empty references", func() {
			var ref Ref[any]

			Convey("Then should have correct properties", func() {
				So(ref.Empty(), ShouldBeTrue)
				So(ref.Type(), ShouldEqual, TypeUnknown)
			})

			Convey("And should return nil for all types", func() {
				So(ref.AsLeaf(), ShouldBeNil)
				So(ref.AsNode4(), ShouldBeNil)
				So(ref.AsNode16(), ShouldBeNil)
				So(ref.AsNode48(), ShouldBeNil)
				So(ref.AsNode256(), ShouldBeNil)
				So(ref.AsNode(), ShouldBeNil)
			})

			Convey("And should not be any specific type", func() {
				So(ref.IsLeaf(), ShouldBeFalse)
				So(ref.IsNode4(), ShouldBeFalse)
				So(ref.IsNode16(), ShouldBeFalse)
				So(ref.IsNode48(), ShouldBeFalse)
				So(ref.IsNode256(), ShouldBeFalse)
				So(ref.IsNode(), ShouldBeFalse)
			})
		})

		Convey("When using AsNode method", func() {
			Convey("And reference is a leaf", func() {
				leaf := NewLeaf[any](a, []byte("hello"), 123)
				ref := NewRef[any](TypeLeaf, leaf)

				nodeResult := ref.AsNode()
				So(nodeResult, ShouldNotBeNil)
				So(nodeResult.Type(), ShouldEqual, TypeLeaf)
			})

			Convey("And reference is a Node4", func() {
				node4 := arena.New(a, Node4[any]{})
				ref := NewRef[any](TypeNode4, node4)

				nodeResult := ref.AsNode()
				So(nodeResult, ShouldNotBeNil)
				So(nodeResult.Type(), ShouldEqual, TypeNode4)
			})

			Convey("And reference is a Node16", func() {
				node16 := arena.New(a, Node16[any]{})
				ref := NewRef[any](TypeNode16, node16)

				nodeResult := ref.AsNode()
				So(nodeResult, ShouldNotBeNil)
				So(nodeResult.Type(), ShouldEqual, TypeNode16)
			})

			Convey("And reference is a Node48", func() {
				node48 := arena.New(a, Node48[any]{})
				ref := NewRef[any](TypeNode48, node48)

				nodeResult := ref.AsNode()
				So(nodeResult, ShouldNotBeNil)
				So(nodeResult.Type(), ShouldEqual, TypeNode48)
			})

			Convey("And reference is a Node256", func() {
				node256 := arena.New(a, Node256[any]{})
				ref := NewRef[any](TypeNode256, node256)

				nodeResult := ref.AsNode()
				So(nodeResult, ShouldNotBeNil)
				So(nodeResult.Type(), ShouldEqual, TypeNode256)
			})

			Convey("And reference is empty", func() {
				var ref Ref[any]
				nodeResult := ref.AsNode()
				So(nodeResult, ShouldBeNil)
			})
		})
	})
}

func TestRef_Replace(t *testing.T) {
	Convey("Given Ref Replace method", t, func() {
		a := &arena.Arena{}

		Convey("When replacing with new reference", func() {
			Convey("And replacing leaf with node", func() {
				leaf := NewLeaf[any](a, []byte("hello"), 123)
				ref := leaf.Ref()

				node4 := arena.New(a, Node4[any]{})
				oldNode := ref.Replace(node4)

				Convey("Then should return old node", func() {
					So(oldNode, ShouldEqual, leaf)
				})

				Convey("And reference should be updated", func() {
					So(ref.Type(), ShouldEqual, TypeNode4)
					So(ref.IsNode4(), ShouldBeTrue)
					So(ref.AsNode4(), ShouldEqual, node4)
				})
			})

			Convey("And replacing node with leaf", func() {
				node4 := arena.New(a, Node4[any]{})
				ref := node4.Ref()

				leaf := NewLeaf[any](a, []byte("world"), 456)
				oldNode := ref.Replace(leaf)

				Convey("Then should return old node", func() {
					So(oldNode, ShouldEqual, node4)
				})

				Convey("And reference should be updated", func() {
					So(ref.Type(), ShouldEqual, TypeLeaf)
					So(ref.IsLeaf(), ShouldBeTrue)
					So(ref.AsLeaf(), ShouldEqual, leaf)
				})
			})

			Convey("And replacing with nil", func() {
				leaf := NewLeaf[any](a, []byte("hello"), 123)
				ref := leaf.Ref()

				oldNode := ref.Replace(nil)

				Convey("Then should return old node", func() {
					So(oldNode, ShouldEqual, leaf)
				})

				Convey("And reference should be empty", func() {
					So(ref.Empty(), ShouldBeTrue)
					So(ref.Type(), ShouldEqual, TypeUnknown)
				})
			})

			Convey("And replacing empty reference", func() {
				var ref Ref[any]
				leaf := NewLeaf[any](a, []byte("hello"), 123)

				oldNode := ref.Replace(leaf)

				Convey("Then should return nil", func() {
					So(oldNode, ShouldBeNil)
				})

				Convey("And reference should be updated", func() {
					So(ref.Type(), ShouldEqual, TypeLeaf)
					So(ref.IsLeaf(), ShouldBeTrue)
					So(ref.AsLeaf(), ShouldEqual, leaf)
				})
			})
		})

		Convey("When replacing with same type", func() {
			Convey("And replacing Node4 with another Node4", func() {
				node4a := arena.New(a, Node4[any]{})
				ref := node4a.Ref()

				node4b := arena.New(a, Node4[any]{})
				oldNode := ref.Replace(node4b)

				Convey("Then should return old node", func() {
					So(oldNode, ShouldEqual, node4a)
				})

				Convey("And reference should be updated", func() {
					So(ref.Type(), ShouldEqual, TypeNode4)
					So(ref.AsNode4(), ShouldEqual, node4b)
				})
			})

			Convey("And replacing Node16 with another Node16", func() {
				node16a := arena.New(a, Node16[any]{})
				ref := node16a.Ref()

				node16b := arena.New(a, Node16[any]{})
				oldNode := ref.Replace(node16b)

				Convey("Then should return old node", func() {
					So(oldNode, ShouldEqual, node16a)
				})

				Convey("And reference should be updated", func() {
					So(ref.Type(), ShouldEqual, TypeNode16)
					So(ref.AsNode16(), ShouldEqual, node16b)
				})
			})
		})
	})
}

func TestRef_EdgeCases(t *testing.T) {
	Convey("Given Ref edge cases", t, func() {
		a := &arena.Arena{}

		Convey("When working with invalid node types", func() {
			Convey("And creating reference with invalid type", func() {
				leaf := NewLeaf(a, []byte("hello"), 123)
				// Create a reference with an invalid type by manipulating the bits
				invalidRef := Ref[any](uintptr(leaf.Ref()) | 0x7) // Invalid type bits

				Convey("Then AsNode should panic", func() {
					So(func() { invalidRef.AsNode() }, ShouldPanicWith, "invalid node type")
				})
			})
		})

		Convey("When working with very large addresses", func() {
			Convey("And creating reference with large address", func() {
				// This test verifies that the address masking works correctly
				leaf := NewLeaf(a, []byte("hello"), 123)
				ref := leaf.Ref()

				Convey("Then type should be preserved", func() {
					So(ref.Type(), ShouldEqual, TypeLeaf)
					So(ref.IsLeaf(), ShouldBeTrue)
				})

				Convey("And node should be accessible", func() {
					leafResult := ref.AsLeaf()
					So(leafResult, ShouldEqual, leaf)
				})
			})
		})

		Convey("When working with zero addresses", func() {
			Convey("And creating reference with zero address", func() {
				// This test verifies that zero addresses are handled correctly
				zeroRef := Ref[any](0)

				Convey("Then should be empty", func() {
					So(zeroRef.Empty(), ShouldBeTrue)
					So(zeroRef.Type(), ShouldEqual, TypeUnknown)
				})

				Convey("And all type checks should return false", func() {
					So(zeroRef.IsLeaf(), ShouldBeFalse)
					So(zeroRef.IsNode4(), ShouldBeFalse)
					So(zeroRef.IsNode16(), ShouldBeFalse)
					So(zeroRef.IsNode48(), ShouldBeFalse)
					So(zeroRef.IsNode256(), ShouldBeFalse)
					So(zeroRef.IsNode(), ShouldBeFalse)
				})

				Convey("And all accessors should return nil", func() {
					So(zeroRef.AsLeaf(), ShouldBeNil)
					So(zeroRef.AsNode4(), ShouldBeNil)
					So(zeroRef.AsNode16(), ShouldBeNil)
					So(zeroRef.AsNode48(), ShouldBeNil)
					So(zeroRef.AsNode256(), ShouldBeNil)
					So(zeroRef.AsNode(), ShouldBeNil)
				})
			})
		})
	})
}

func TestRef_Performance(t *testing.T) {
	Convey("Given Ref performance considerations", t, func() {
		a := &arena.Arena{}

		Convey("When performing many type checks", func() {
			Convey("And checking type of 1000 references", func() {
				refs := make([]Ref[any], 1000)
				types := []Type{TypeLeaf, TypeNode4, TypeNode16, TypeNode48, TypeNode256}

				// Create references of different types
				for i := 0; i < 1000; i++ {
					switch types[i%len(types)] {
					case TypeLeaf:
						leaf := NewLeaf[any](a, []byte{byte(i)}, i)
						refs[i] = leaf.Ref()
					case TypeNode4:
						node4 := arena.New(a, Node4[any]{})
						refs[i] = node4.Ref()
					case TypeNode16:
						node16 := arena.New(a, Node16[any]{})
						refs[i] = node16.Ref()
					case TypeNode48:
						node48 := arena.New(a, Node48[any]{})
						refs[i] = node48.Ref()
					case TypeNode256:
						node256 := arena.New(a, Node256[any]{})
						refs[i] = node256.Ref()
					}
				}

				So(len(refs), ShouldEqual, 1000)

				// Perform type checks
				leafCount := 0
				node4Count := 0
				node16Count := 0
				node48Count := 0
				node256Count := 0

				for _, ref := range refs {
					if ref.IsLeaf() {
						leafCount++
					} else if ref.IsNode4() {
						node4Count++
					} else if ref.IsNode16() {
						node16Count++
					} else if ref.IsNode48() {
						node48Count++
					} else if ref.IsNode256() {
						node256Count++
					}
				}

				// Verify distribution (should be roughly 200 of each type)
				So(leafCount, ShouldBeGreaterThan, 150)
				So(node4Count, ShouldBeGreaterThan, 150)
				So(node16Count, ShouldBeGreaterThan, 150)
				So(node48Count, ShouldBeGreaterThan, 150)
				So(node256Count, ShouldBeGreaterThan, 150)
			})

			Convey("And performing many replacements", func() {
				leaf := NewLeaf[any](a, []byte("hello"), 123)
				ref := leaf.Ref()

				// Perform many replacements
				for i := 0; i < 100; i++ {
					if i%2 == 0 {
						node4 := arena.New(a, Node4[any]{})
						ref.Replace(node4)
					} else {
						newLeaf := NewLeaf[any](a, []byte{byte(i)}, i)
						ref.Replace(newLeaf)
					}
				}

				Convey("Then reference should still be valid", func() {
					So(ref.Empty(), ShouldBeFalse)
					// The final type depends on the last replacement
					So(ref.Type() == TypeLeaf || ref.Type() == TypeNode4, ShouldBeTrue)
				})
			})
		})
	})
}
