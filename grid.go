package waveform

import (
	"image"
	"image/color"
)

type Stack struct {
	width, height int
	imgs          []image.Image
}

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

func (s *Stack) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{s.width, s.height},
	}
}

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

func (s *Stack) ColorModel() color.Model {
	return color.RGBAModel
}
