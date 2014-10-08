package waveform

import (
	"image"
	"image/color"
	"sync"
)

var (
	DefaultFillColor = color.RGBAModel.Convert(color.Black)
	DefaultBGColor   = color.RGBAModel.Convert(color.White)
)

type Pair struct {
	Max, Min float32
}

type WaveForm struct {
	// Protects writes, this can be locked by user to avoid
	// a write happening while reading out the entire image
	sync.Mutex

	FillColor       color.Color
	BackgroundColor color.Color
	Width, Height   int

	// img is a slice of min/max pairs
	img []Pair
	// state to keep track of for AddPair

	// startPos is the x the image starts in
	startPos int
	// full determines if img is full and we need to wrap around
	full bool
	// imgPos is the next position to write to in img
	imgPos int
}

// NewWaveForm returns an initialized empty WaveForm image
//
// The width and height are the dimensions of the image returned.
func NewWaveForm(width, height int) *WaveForm {
	return &WaveForm{
		FillColor:       DefaultFillColor,
		BackgroundColor: DefaultBGColor,
		img:             make([]Pair, width),
		Width:           width,
		Height:          height,
	}
}

func (wf *WaveForm) Resize(width, height int) {
	wf.Lock()
	defer wf.Unlock()

	// We only need to do actual work if the width changes
	if width == wf.Width {
		wf.Height = height
		return
	}

	nimg := make([]Pair, width)

	copy(nimg, wf.img[wf.startPos:])
	copy(nimg[len(wf.img)-wf.startPos:], wf.img[:wf.startPos])

	wf.img = nimg
	wf.Width = width
	wf.Height = height
}

// ColorModel implements the image.Image ColorModel method.
func (wf *WaveForm) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds implements the image.Image Bounds method.
func (wf *WaveForm) Bounds() image.Rectangle {
	return image.Rectangle{
		image.Point{0, 0},
		image.Point{int(wf.Width), int(wf.Height)},
	}
}

// At implements the image.Image At method.
func (wf *WaveForm) At(x, y int) color.Color {
	x = x + wf.startPos
	if x > wf.Width {
		x = x % wf.Width
	}
	p := wf.img[x]
	mid := wf.Height / 2

	if y == mid {
		return wf.FillColor
	}

	if y < mid && mid+int(float32(mid)*p.Min) < y {
		return wf.FillColor
	}

	if y > mid && mid+int(float32(mid)*p.Max) > y {
		return wf.FillColor
	}

	return wf.BackgroundColor
}

// AddPair adds a pair of min and max points, each pair is
// a single pixel in the resulting image.
//
// AddPair can be called wf.width times before it gets rid
// of previous pairs added.
func (wf *WaveForm) AddPair(p Pair) {
	wf.Lock()
	x := wf.imgPos
	wf.imgPos++

	// wrap around our x if needed
	if wf.imgPos >= len(wf.img) {
		wf.full = true
		wf.imgPos = 0
	}

	// if our buffer is already full we will need to adjust things
	if wf.full {
		wf.startPos = wf.imgPos
	}

	wf.img[x] = p
	wf.Unlock()
}
