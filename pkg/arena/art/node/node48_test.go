package node

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestNode48(t *testing.T) {
	Convey("Given a Node48", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node48{})

		Convey("When checking basic properties", func() {
			So(node.Type(), ShouldEqual, TypeNode48)
			So(node.Full(), ShouldBeFalse)
			So(node.NumChildren, ShouldEqual, 0)
			So(node.Ref().Type(), ShouldEqual, TypeNode48)
		})

		Convey("When adding children", func() {
			// Create mock children
			children := make([]*Leaf, 48)
			for i := 0; i < 48; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte(i))})
			}

			Convey("Adding first child", func() {
				node.AddChild(42, children[0])
				So(node.NumChildren, ShouldEqual, 1)
				So(node.Keys[42], ShouldEqual, byte(1)) // 1-based indexing
				So(node.Children[0], ShouldEqual, children[0].Ref())
			})

			Convey("Adding multiple children", func() {
				node.AddChild(10, children[0])
				node.AddChild(20, children[1])
				node.AddChild(30, children[2])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Keys[10], ShouldEqual, byte(1))
				So(node.Keys[20], ShouldEqual, byte(2))
				So(node.Keys[30], ShouldEqual, byte(3))
				So(node.Children[0], ShouldEqual, children[0].Ref())
				So(node.Children[1], ShouldEqual, children[1].Ref())
				So(node.Children[2], ShouldEqual, children[2].Ref())
			})

			Convey("Adding children with sparse keys", func() {
				// Add children with widely spaced keys
				node.AddChild(0, children[0])
				node.AddChild(128, children[1])
				node.AddChild(255, children[2])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Keys[0], ShouldEqual, byte(1))
				So(node.Keys[128], ShouldEqual, byte(2))
				So(node.Keys[255], ShouldEqual, byte(3))
			})

			Convey("Adding children to fill capacity", func() {
				// Add 48 children to fill the node
				for i := 0; i < 48; i++ {
					node.AddChild(byte(i), children[i])
				}

				So(node.NumChildren, ShouldEqual, 48)
				So(node.Full(), ShouldBeTrue)

				// Verify all children can be found
				for i := 0; i < 48; i++ {
					found := node.FindChild(byte(i))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})
		})

		Convey("When finding children", func() {
			// Setup children
			children := make([]*Leaf, 10)
			for i := 0; i < 10; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte(i*25))})
				node.AddChild(byte(i*25), children[i])
			}

			Convey("Finding existing children", func() {
				for i := 0; i < 10; i++ {
					found := node.FindChild(byte(i * 25))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})

			Convey("Finding non-existent children", func() {
				// Test keys that are not in the sparse distribution
				found := node.FindChild(1)
				So(found, ShouldBeNil)

				found = node.FindChild(26) // Between 0 and 25
				So(found, ShouldBeNil)

				found = node.FindChild(51) // Between 25 and 50
				So(found, ShouldBeNil)
			})

			Convey("Finding children at boundaries", func() {
				// Test finding children at byte boundaries
				So(node.FindChild(0), ShouldNotBeNil)   // First child
				So(node.FindChild(225), ShouldNotBeNil) // Last child
			})
		})

		Convey("When checking capacity", func() {
			Convey("Empty node is not full", func() {
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 24 children is not full", func() {
				for i := 0; i < 24; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte(i*10))})
					node.AddChild(byte(i*10), child)
				}
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 48 children is full", func() {
				for i := 0; i < 48; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte(i*5))})
					node.AddChild(byte(i*5), child)
				}
				So(node.Full(), ShouldBeTrue)
			})
		})

		Convey("When growing to Node256", func() {
			// Setup children
			for i := 0; i < 48; i++ {
				child := arena.New(a, Leaf{Key: slice.Of(a, byte(i*5))})
				node.AddChild(byte(i*5), child)
			}

			Convey("Growing should create Node256", func() {
				newNode := node.Grow(a)
				So(newNode.Type(), ShouldEqual, TypeNode256)
			})

			Convey("Growing should preserve all children", func() {
				newNode := node.Grow(a)
				node256 := newNode.(*Node256)

				So(node256.NumChildren, ShouldEqual, 48)
				// Check that all children are properly mapped
				for i := 0; i < 48; i++ {
					key := byte(i * 5)
					So(node256.Children[key], ShouldEqual, node.Children[i])
				}
			})
		})

		Convey("When getting minimum and maximum", func() {
			Convey("Empty node should return nil", func() {
				So(node.Minimum(), ShouldBeNil)
				So(node.Maximum(), ShouldBeNil)
			})

			Convey("Node with children should return correct min/max", func() {
				child1 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})
				child2 := arena.New(a, Leaf{Key: slice.FromString(a, "b")})
				child3 := arena.New(a, Leaf{Key: slice.FromString(a, "c")})

				node.AddChild('c', child3)
				node.AddChild('a', child1)
				node.AddChild('b', child2)

				So(node.Minimum(), ShouldEqual, child1)
				So(node.Maximum(), ShouldEqual, child3)
			})
		})
	})
}

func TestNode48_EdgeCases(t *testing.T) {
	Convey("Given a Node48 with edge cases", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node48{})

		Convey("When adding duplicate keys", func() {
			child1 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})
			child2 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})

			node.AddChild('a', child1)
			node.AddChild('a', child2)

			// Should replace the existing child
			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild('a')
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child2.Ref())
		})

		Convey("When adding zero byte key", func() {
			child := arena.New(a, Leaf{Key: slice.Of(a, byte(0))})
			node.AddChild(0, child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(0)
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding 255 byte key", func() {
			child := arena.New(a, Leaf{Key: slice.Of(a, byte(255))})
			node.AddChild(255, child)

			So(node.NumChildren, ShouldEqual, 1)
			found := node.FindChild(255)
			So(found, ShouldNotBeNil)
			So(*found, ShouldEqual, child.Ref())
		})

		Convey("When adding children at sparse intervals", func() {
			// Add children with very sparse distribution
			sparseKeys := []byte{0, 64, 128, 192, 255}
			for i, key := range sparseKeys {
				child := arena.New(a, Leaf{Key: slice.Of(a, key)})
				node.AddChild(key, child)
				So(node.NumChildren, ShouldEqual, i+1)
			}

			// Verify all sparse children can be found
			for _, key := range sparseKeys {
				found := node.FindChild(key)
				So(found, ShouldNotBeNil)
			}
		})
	})
}

func TestNode48_Performance(t *testing.T) {
	Convey("Given a Node48 with performance considerations", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node48{})

		Convey("When adding many children", func() {
			// Test that adding 48 children works correctly
			children := make([]*Leaf, 48)
			for i := 0; i < 48; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte(i*5))})
				node.AddChild(byte(i*5), children[i])
			}

			So(node.NumChildren, ShouldEqual, 48)
			So(node.Full(), ShouldBeTrue)

			// Verify all children can be found
			for i := 0; i < 48; i++ {
				found := node.FindChild(byte(i * 5))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, children[i].Ref())
			}
		})

		Convey("When searching with sparse distribution", func() {
			// Add children with sparse distribution to test search performance
			sparseKeys := []byte{1, 3, 7, 15, 31, 63, 127, 255}
			for _, key := range sparseKeys {
				child := arena.New(a, Leaf{Key: slice.Of(a, key)})
				node.AddChild(key, child)
			}

			// Test finding existing and non-existing keys
			for _, key := range sparseKeys {
				So(node.FindChild(key), ShouldNotBeNil)
			}

			// Test some non-existing keys
			So(node.FindChild(2), ShouldBeNil)
			So(node.FindChild(8), ShouldBeNil)
			So(node.FindChild(16), ShouldBeNil)
		})

		Convey("When testing 1-based indexing", func() {
			// Verify that the 1-based indexing works correctly
			child1 := arena.New(a, Leaf{Key: slice.Of(a, byte(100))})
			child2 := arena.New(a, Leaf{Key: slice.Of(a, byte(200))})

			node.AddChild(100, child1)
			node.AddChild(200, child2)

			So(node.Keys[100], ShouldEqual, byte(1))
			So(node.Keys[200], ShouldEqual, byte(2))
			So(node.Children[0], ShouldEqual, child1.Ref())
			So(node.Children[1], ShouldEqual, child2.Ref())
		})
	})
}
