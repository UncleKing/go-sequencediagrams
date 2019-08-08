package utils

import (
	"image"
)

type Rectangle struct {
	image.Rectangle
}

func (r Rectangle) MidX() int {
	return r.Min.X + r.Dx()/2
}

func (r Rectangle) MidY() int {
	return r.Min.Y + r.Dy()/2
}

// Rect is shorthand for Rectangle{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect(x0, y0, x1, y1 int) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	r := Rectangle{}
	r.Min = image.Point{x0, y0}
	r.Max = image.Point{x1, y1}
	return r
}

// Add returns the rectangle r translated by p.
func (r Rectangle) Add(p image.Point) Rectangle {

	s := Rectangle{}
	s.Min = image.Point{X: r.Min.X + p.X, Y: r.Min.Y + p.Y}
	s.Max = image.Point{X: r.Max.X + p.X, Y: r.Max.Y + p.Y}
	return s
}
