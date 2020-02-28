package tree

// Inspired by https://www.azurefromthetrenches.com/introductory-guide-to-aabb-tree-collision-detection/

import (
	"errors"
	"math"
	"reflect"
)

var (
	ErrtNotInTree    = errors.New("The object is not in the tree")
	ErrAlreadyInTree = errors.New("The object is already in the tree")
)

type AABB interface {
	AABB() *AABBData
}

type AABBData struct {
	MinX, MinY, MaxX, MaxY float64
}

func (aabb *AABBData) AABB() *AABBData {
	return aabb
}
func (aabb *AABBData) SurfaceArea() float64 {
	a := aabb.MaxX - aabb.MinX
	b := aabb.MaxY - aabb.MinY
	return a * b
}

func (aabb *AABBData) Merge(other *AABBData) *AABBData {
	return Merge(aabb, other)
}

func Merge(a, b AABB) *AABBData {
	aData := a.AABB()
	bData := b.AABB()
	return &AABBData{
		MinX: math.Min(aData.MinX, bData.MinX),
		MinY: math.Min(aData.MinY, bData.MinY),
		MaxX: math.Max(aData.MaxX, bData.MaxX),
		MaxY: math.Max(aData.MaxY, bData.MaxY),
	}
}

func (aabb *AABBData) Overlaps(other *AABBData) bool {
	return Overlaps(aabb, other)
}

func Overlaps(a, b AABB) bool {
	aData := a.AABB()
	bData := b.AABB()
	return aData.MaxX > bData.MinX &&
		aData.MinX < bData.MaxX &&
		aData.MaxY > bData.MinY &&
		aData.MinY < bData.MaxY
}

func (aabb *AABBData) Move(x, y float64) *AABBData {
	return Move(aabb, x, y)
}

func Move(obj AABB, x, y float64) *AABBData {
	data := obj.AABB()
	return &AABBData{
		MinX: data.MinX + x,
		MaxX: data.MaxX + x,
		MinY: data.MinY + y,
		MaxY: data.MaxY + y,
	}
}

type aabbTreeNode struct {
	Object     AABB      `json:"-"`
	ObjectAABB *AABBData `json:"aabb"`

	Parent *aabbTreeNode `json:"-"`
	Left   *aabbTreeNode `json:"left"`
	Right  *aabbTreeNode `json:"right"`

	Depth int `json:"depth"`
}

func newAABBTreeNode(object AABB) *aabbTreeNode {
	if reflect.ValueOf(object).Kind() != reflect.Ptr {
		panic("provided object must be a pointer")
	}
	return &aabbTreeNode{
		Object:     object,
		ObjectAABB: object.AABB(),
	}
}

func (node *aabbTreeNode) IsLeaf() bool {
	return node.Left == nil
}
func (node *aabbTreeNode) AABB() *AABBData {
	return node.ObjectAABB
}
