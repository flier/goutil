package node

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestNode16(t *testing.T) {
	Convey("Given a Node16", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node16{})

		Convey("When checking basic properties", func() {
			So(node.Type(), ShouldEqual, TypeNode16)
			So(node.Full(), ShouldBeFalse)
			So(node.NumChildren, ShouldEqual, 0)
			So(node.Ref().Type(), ShouldEqual, TypeNode16)
		})

		Convey("When adding children", func() {
			// Create mock children
			children := make([]*Leaf, 16)
			for i := 0; i < 16; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte('a'+i))})
			}

			Convey("Adding first child", func() {
				node.AddChild('a', children[0])
				So(node.NumChildren, ShouldEqual, 1)
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Children[0], ShouldEqual, children[0].Ref())
			})

			Convey("Adding children in order", func() {
				for i := 0; i < 8; i++ {
					node.AddChild(byte('a'+i), children[i])
				}

				So(node.NumChildren, ShouldEqual, 8)
				for i := 0; i < 8; i++ {
					So(node.Keys[i], ShouldEqual, byte('a'+i))
					So(node.Children[i], ShouldEqual, children[i].Ref())
				}
			})

			Convey("Adding children out of order", func() {
				node.AddChild('c', children[2])
				node.AddChild('a', children[0])
				node.AddChild('b', children[1])

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Keys[1], ShouldEqual, byte('b'))
				So(node.Keys[2], ShouldEqual, byte('c'))
			})

			Convey("Adding children to maintain sorted order", func() {
				// Add in reverse order
				for i := 15; i >= 0; i-- {
					node.AddChild(byte('a'+i), children[i])
				}

				So(node.NumChildren, ShouldEqual, 16)
				for i := 0; i < 16; i++ {
					So(node.Keys[i], ShouldEqual, byte('a'+i))
					So(node.Children[i], ShouldEqual, children[i].Ref())
				}
			})
		})

		Convey("When finding children", func() {
			// Setup children
			children := make([]*Leaf, 8)
			for i := 0; i < 8; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte('a'+i))})
				node.AddChild(byte('a'+i), children[i])
			}

			Convey("Finding existing children", func() {
				for i := 0; i < 8; i++ {
					found := node.FindChild(byte('a' + i))
					So(found, ShouldNotBeNil)
					So(*found, ShouldEqual, children[i].Ref())
				}
			})

			Convey("Finding non-existent children", func() {
				found := node.FindChild('x')
				So(found, ShouldBeNil)

				found = node.FindChild('z')
				So(found, ShouldBeNil)
			})
		})

		Convey("When checking capacity", func() {
			Convey("Empty node is not full", func() {
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 8 children is not full", func() {
				for i := 0; i < 8; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte('a'+i))})
					node.AddChild(byte('a'+i), child)
				}
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 16 children is full", func() {
				for i := 0; i < 16; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte('a'+i))})
					node.AddChild(byte('a'+i), child)
				}
				So(node.Full(), ShouldBeTrue)
			})
		})

		Convey("When growing to Node48", func() {
			// Setup children
			for i := 0; i < 16; i++ {
				child := arena.New(a, Leaf{Key: slice.Of(a, byte('a'+i))})
				node.AddChild(byte('a'+i), child)
			}

			Convey("Growing should create Node48", func() {
				newNode := node.Grow(a)
				So(newNode.Type(), ShouldEqual, TypeNode48)
			})

			Convey("Growing should preserve all children", func() {
				newNode := node.Grow(a)
				node48 := newNode.(*Node48)

				So(node48.NumChildren, ShouldEqual, 16)
				// Check that all children are properly mapped
				for i := 0; i < 16; i++ {
					key := byte('a' + i)
					So(node48.Keys[key], ShouldNotEqual, 0)
					childIndex := node48.Keys[key] - 1
					So(node48.Children[childIndex], ShouldEqual, node.Children[i])
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

func TestNode16_EdgeCases(t *testing.T) {
	Convey("Given a Node16 with edge cases", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node16{})

		Convey("When adding duplicate keys", func() {
			child1 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})
			child2 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})

			node.AddChild('a', child1)
			node.AddChild('a', child2)

			Convey("Should not replace the existing child", func() {
				So(node.NumChildren, ShouldEqual, 2)

				found := node.FindChild('a')
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child1.Ref())
			})
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

		Convey("When adding children at boundaries", func() {
			// Add children at the beginning and end of the byte range
			childStart := arena.New(a, Leaf{Key: slice.Of(a, byte(0))})
			childEnd := arena.New(a, Leaf{Key: slice.Of(a, byte(255))})

			node.AddChild(0, childStart)
			node.AddChild(255, childEnd)

			So(node.NumChildren, ShouldEqual, 2)
			So(node.Keys[0], ShouldEqual, byte(0))
			So(node.Keys[1], ShouldEqual, byte(255))
		})
	})
}

func TestNode16_Performance(t *testing.T) {
	Convey("Given a Node16 with performance considerations", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node16{})

		Convey("When adding many children", func() {
			// Test that adding 16 children works correctly
			children := make([]*Leaf, 16)
			for i := 0; i < 16; i++ {
				children[i] = arena.New(a, Leaf{Key: slice.Of(a, byte(i))})
				node.AddChild(byte(i), children[i])
			}

			So(node.NumChildren, ShouldEqual, 16)
			So(node.Full(), ShouldBeTrue)

			// Verify all children can be found
			for i := 0; i < 16; i++ {
				found := node.FindChild(byte(i))
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, children[i].Ref())
			}
		})

		Convey("When searching in sorted array", func() {
			// Add children in sorted order to test search performance
			for i := 0; i < 8; i++ {
				child := arena.New(a, Leaf{Key: slice.Of(a, byte(i*2))})
				node.AddChild(byte(i*2), child)
			}

			// Test finding existing and non-existing keys
			So(node.FindChild(0), ShouldNotBeNil)
			So(node.FindChild(14), ShouldNotBeNil)
			So(node.FindChild(1), ShouldBeNil)  // Odd numbers don't exist
			So(node.FindChild(15), ShouldBeNil) // Odd numbers don't exist
		})
	})
}
