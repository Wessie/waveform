package waveform

import (
	"image"
	"image/color"
)

// Stack is a simple image.Image implementation that stacks other
// image.Image's on top of each other vertically.
//
// The width of the Stack is of the widest image given. The height
// of the Stack is the height of each image added up.
//
// Points that fall outside of the image.Image return an undefined
// color.Color, therefore it is suggested to only use images of the
// same width.
type Stack struct {
	width, height int
	imgs          []image.Image
}

// NewStack creates a new Stack from the images passed in.
func NewStack(imgs ...image.Image) image.Image {
	var w, h int
	for i := range imgs {
		h += imgs[i].Bounds().Max.Y
		x := imgs[i].Bounds().Max.X
		if x > w {
			w = x
		}
	}

	return &Stack{
		imgs:   imgs,
		width:  w,
		height: h,
	}
}

// Bounds implements the image.Image method
func (s *Stack) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{s.width, s.height},
	}
}

// At implements the image.Image method
func (s *Stack) At(x, y int) color.Color {
	var passedY int
	for i := range s.imgs {
		b := s.imgs[i].Bounds().Max.Y + passedY

		if y > b {
			passedY = b
			continue
		}

		y = y - passedY

		return s.imgs[i].At(x, y)
	}

	return nil
}

// ColorModel implements the image.Image method
func (s *Stack) ColorModel() color.Model {
	return color.RGBAModel
}
