package aabb

// Inspired by https://www.azurefromthetrenches.com/introductory-guide-to-aabb-tree-collision-detection/

import (
	"errors"
	"math"
)

var (
	ErrtNotInTree    = errors.New("The object is not in the tree")
	ErrAlreadyInTree = errors.New("The object is already in the tree")
	ErrNotAReference = errors.New("Objects must be passed by reference")
)

type AABB interface {
	AABB() *AABBData
}

// AABBData represents the AABB coordinates.
type AABBData struct {
	MinX, MinY, MaxX, MaxY float64
}

func NewAABBData(MinX, MinY, MaxX, MaxY float64) (*AABBData, error) {
	data := &AABBData{MinX, MinY, MaxX, MaxY}
	if !data.IsValid() {
		return nil, errors.New("Shape is invalid, did you swap Min and Max values?")
	}
	return data, nil
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

func (aabb *AABBData) IsValid() bool {
	return aabb.MaxX > aabb.MinY && aabb.MaxY > aabb.MinY
}
