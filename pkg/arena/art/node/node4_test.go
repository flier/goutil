package node

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/flier/goutil/pkg/arena"
	"github.com/flier/goutil/pkg/arena/slice"
)

func TestNode4(t *testing.T) {
	Convey("Given a Node4", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node4{})

		Convey("When checking basic properties", func() {
			So(node.Type(), ShouldEqual, TypeNode4)
			So(node.Full(), ShouldBeFalse)
			So(node.NumChildren, ShouldEqual, 0)
			So(node.Ref().Type(), ShouldEqual, TypeNode4)
		})

		Convey("When adding children", func() {
			// Create mock children
			child1 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})
			child2 := arena.New(a, Leaf{Key: slice.FromString(a, "b")})
			child3 := arena.New(a, Leaf{Key: slice.FromString(a, "c")})
			child4 := arena.New(a, Leaf{Key: slice.FromString(a, "d")})

			Convey("Adding first child", func() {
				node.AddChild('a', child1)
				So(node.NumChildren, ShouldEqual, 1)
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Children[0], ShouldEqual, child1.Ref())
			})

			Convey("Adding children in order", func() {
				node.AddChild('a', child1)
				node.AddChild('b', child2)
				node.AddChild('c', child3)

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Keys[1], ShouldEqual, byte('b'))
				So(node.Keys[2], ShouldEqual, byte('c'))
			})

			Convey("Adding children out of order", func() {
				node.AddChild('c', child3)
				node.AddChild('a', child1)
				node.AddChild('b', child2)

				So(node.NumChildren, ShouldEqual, 3)
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Keys[1], ShouldEqual, byte('b'))
				So(node.Keys[2], ShouldEqual, byte('c'))
			})

			Convey("Adding children to maintain sorted order", func() {
				node.AddChild('d', child4)
				node.AddChild('b', child2)
				node.AddChild('a', child1)
				node.AddChild('c', child3)

				So(node.NumChildren, ShouldEqual, 4)
				So(node.Keys[0], ShouldEqual, byte('a'))
				So(node.Keys[1], ShouldEqual, byte('b'))
				So(node.Keys[2], ShouldEqual, byte('c'))
				So(node.Keys[3], ShouldEqual, byte('d'))
			})
		})

		Convey("When finding children", func() {
			// Setup children
			child1 := arena.New(a, Leaf{Key: slice.FromString(a, "a")})
			child2 := arena.New(a, Leaf{Key: slice.FromString(a, "b")})
			child3 := arena.New(a, Leaf{Key: slice.FromString(a, "c")})

			node.AddChild('a', child1)
			node.AddChild('b', child2)
			node.AddChild('c', child3)

			Convey("Finding existing children", func() {
				found := node.FindChild('a')
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child1.Ref())

				found = node.FindChild('b')
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child2.Ref())

				found = node.FindChild('c')
				So(found, ShouldNotBeNil)
				So(*found, ShouldEqual, child3.Ref())
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

			Convey("Node with 3 children is not full", func() {
				for i := 0; i < 3; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte('a'+i))})
					node.AddChild(byte('a'+i), child)
				}
				So(node.Full(), ShouldBeFalse)
			})

			Convey("Node with 4 children is full", func() {
				for i := 0; i < 4; i++ {
					child := arena.New(a, Leaf{Key: slice.Of(a, byte('a'+i))})
					node.AddChild(byte('a'+i), child)
				}
				So(node.Full(), ShouldBeTrue)
			})
		})

		Convey("When growing to Node16", func() {
			// Setup children
			for i := 0; i < 4; i++ {
				child := arena.New(a, Leaf{Key: slice.Of(a, byte('a'+i))})
				node.AddChild(byte('a'+i), child)
			}

			Convey("Growing should create Node16", func() {
				newNode := node.Grow(a)
				So(newNode.Type(), ShouldEqual, TypeNode16)
			})

			Convey("Growing should preserve all children", func() {
				newNode := node.Grow(a)
				node16 := newNode.(*Node16)

				So(node16.NumChildren, ShouldEqual, 4)
				So(node16.Keys[0], ShouldEqual, byte('a'))
				So(node16.Keys[1], ShouldEqual, byte('b'))
				So(node16.Keys[2], ShouldEqual, byte('c'))
				So(node16.Keys[3], ShouldEqual, byte('d'))
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

				// Note: Since these are mock leaves, we can't easily test the actual
				// minimum/maximum values without more complex setup
				So(node.Minimum(), ShouldEqual, child1)
				So(node.Maximum(), ShouldEqual, child3)
			})
		})
	})
}

func TestNode4_EdgeCases(t *testing.T) {
	Convey("Given a Node4 with edge cases", t, func() {
		a := &arena.Arena{}
		node := arena.New(a, Node4{})

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
	})
}
