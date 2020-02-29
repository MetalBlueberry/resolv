package resolv

import (
	"math"
)

// Resolve attempts to move the checking Shape with the specified X and Y values, returning a Collision object
// if it collides with the specified other Shape. The deltaX and deltaY arguments are the movement displacement
// in pixels. For platformers in particular, you would probably want to resolve on the X and Y axes separately.
func Resolve(firstShape Shape, other Shape, deltaX, deltaY float64) Collision {

	out := Collision{}
	out.ResolveX = deltaX
	out.ResolveY = deltaY
	out.ShapeA = firstShape

	if deltaX == 0 && deltaY == 0 {
		return out
	}

	x := deltaX
	y := deltaY

	primeX := true
	slope := 0.0

	if math.Abs(float64(deltaY)) > math.Abs(float64(deltaX)) {
		primeX = false
		if deltaY != 0 && deltaX != 0 {
			slope = deltaX / deltaY
		}
	} else if deltaY != 0 && deltaX != 0 {
		slope = deltaY / deltaX
	}

	for true {

		if firstShape.WouldBeColliding(other, out.ResolveX, out.ResolveY) {

			if primeX {

				if deltaX > 0 {
					x--
				} else if deltaX < 0 {
					x++
				}

				if deltaY > 0 {
					y -= slope
				} else if deltaY < 0 {
					y += slope
				}

			} else {

				if deltaY > 0 {
					y--
				} else if deltaY < 0 {
					y++
				}

				if deltaX > 0 {
					x -= slope
				} else if deltaX < 0 {
					x += slope
				}

			}

			out.ResolveX = x
			out.ResolveY = y
			out.ShapeB = other

		} else {
			break
		}

	}

	if math.Abs(float64(deltaX-out.ResolveX)) > math.Abs(float64(deltaX)*1.5) || math.Abs(float64(deltaY-out.ResolveY)) > math.Abs(float64(deltaY)*1.5) {
		out.Teleporting = true
	}

	return out

}

// Distance returns the distance from one pair of X and Y values to another.
func Distance(x, y, x2, y2 float64) float64 {

	dx := x - x2
	dy := y - y2
	ds := (dx * dx) + (dy * dy)
	return math.Sqrt(math.Abs(float64(ds)))

}
